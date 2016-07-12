package stores

import (
	"sync"

	"github.com/davelondon/kerr"
	"golang.org/x/net/context"
	"kego.io/editor/client/actions"
	"kego.io/editor/client/models"
	"kego.io/flux"
	"kego.io/json"
	"kego.io/system"
	"kego.io/system/node"
)

type NodeStore struct {
	*flux.Store
	ctx context.Context
	app *App

	selected *node.Node
	m        *sync.Mutex
	models   map[*node.Node]*models.NodeModel
}

func (s *NodeStore) Selected() *node.Node {
	return s.selected
}

func (s *NodeStore) Get(n *node.Node) *models.NodeModel {
	s.m.Lock()
	defer s.m.Unlock()
	m, ok := s.models[n]
	if !ok {
		m = &models.NodeModel{Node: n}
		s.models[n] = m
	}
	return m
}

type nodeNotif string

func (b nodeNotif) IsNotif() {}

const (
	NodeInitialised            nodeNotif = "NodeInitialised"
	NodeDeleted                nodeNotif = "NodeDeleted"
	NodeValueChanged           nodeNotif = "NodeValueChanged"
	NodeDescendantValueChanged nodeNotif = "NodeDescendantValueChanged"
	NodeFocus                  nodeNotif = "NodeFocus"
	NodeErrorsChanged          nodeNotif = "NodeErrorsChanged"
	NodeArrayReorder           nodeNotif = "NodeArrayReorder"
	NodeChildAdded             nodeNotif = "NodeChildAdded"
	NodeChildDeleted           nodeNotif = "NodeChildDeleted"
)

func NewNodeStore(ctx context.Context) *NodeStore {
	s := &NodeStore{
		Store:  &flux.Store{},
		ctx:    ctx,
		app:    FromContext(ctx),
		m:      &sync.Mutex{},
		models: map[*node.Node]*models.NodeModel{},
	}
	s.Init(s)
	return s
}

func (s *NodeStore) Handle(payload *flux.Payload) bool {
	switch action := payload.Action.(type) {
	case *actions.Add:
		payload.Wait(s.app.Actions)
		switch action.Direction() {
		case actions.New:
			if err := mutateAddNode(s.ctx, action.Node, action.Parent, action.Key, action.Index, action.Type); err != nil {
				s.app.Fail <- kerr.Wrap("HUOGBUQCAO", err)
				break
			}
			payload.Notify(action.Node, NodeInitialised)
			payload.Notify(action.Parent, NodeChildAdded)
		case actions.Undo:
			action.Backup = node.NewNode()
			if err := mutateDeleteNode(s.ctx, action.Node, action.Parent, action.Backup); err != nil {
				s.app.Fail <- kerr.Wrap("RTAGMUIKMD", err)
				break
			}
			payload.Notify(action.Node, NodeDeleted)
			payload.Notify(action.Parent, NodeChildDeleted)
		case actions.Redo:
			if err := mutateRestoreNode(s.ctx, action.Node, action.Parent, action.Backup); err != nil {
				s.app.Fail <- kerr.Wrap("MHUTMXOGBP", err)
				break
			}
			payload.Notify(action.Node, NodeInitialised)
			payload.Notify(action.Parent, NodeChildAdded)
		}
	case *actions.Delete:
		payload.Wait(s.app.Actions)
		switch action.Direction() {
		case actions.New, actions.Redo:
			action.Backup = node.NewNode()
			if err := mutateDeleteNode(s.ctx, action.Node, action.Parent, action.Backup); err != nil {
				s.app.Fail <- kerr.Wrap("DFHTKJRLQC", err)
				break
			}
			payload.Notify(action.Node, NodeDeleted)
			payload.Notify(action.Parent, NodeChildDeleted)
		case actions.Undo:
			if err := mutateRestoreNode(s.ctx, action.Node, action.Parent, action.Backup); err != nil {
				s.app.Fail <- kerr.Wrap("HAPWUOPBTW", err)
				break
			}
			payload.Notify(action.Node, NodeInitialised)
			payload.Notify(action.Parent, NodeChildAdded)
		}
	case *actions.Reorder:
		payload.Wait(s.app.Actions)
		if !action.Model.Node.Type.IsNativeArray() {
			s.app.Fail <- kerr.New("EPBQVIICFM", "Must be array")
			break
		}
		from := action.Before
		to := action.After
		if action.Direction() == actions.Undo {
			from = action.After
			to = action.Before
		}
		if err := action.Model.Node.ReorderArrayChild(from, to); err != nil {
			s.app.Fail <- kerr.Wrap("HMIFPKVJCN", err)
			break
		}
		payload.Notify(action.Model.Node, NodeArrayReorder)

	case *actions.Modify:
		payload.Wait(s.app.Actions)
		n := action.Editor.Node
		val := action.After
		if action.Direction() == actions.Undo {
			val = action.Before
		}
		if err := n.SetValue(s.ctx, val); err != nil {
			s.app.Fail <- kerr.Wrap("VIMXVIHPFY", err)
			break
		}

		model := s.app.Nodes.Get(n)
		changed, err := model.Validate(s.ctx, s.app.Rule.Get(n.Root(), n))
		if err != nil {
			s.app.Fail <- kerr.Wrap("EEIYMGQCCA", err)
			break
		}
		if changed {
			payload.Notify(n, NodeErrorsChanged)
		}

		payload.Notify(n, NodeValueChanged)

		c := n.Parent
		for c != nil {
			payload.Notify(c, NodeDescendantValueChanged)
			c = c.Parent
		}
	case *actions.BranchSelecting:
		if ni, ok := action.Branch.Contents.(models.NodeContentsInterface); ok {
			s.selected = ni.GetNode()
		} else {
			s.selected = nil
		}
	case *actions.BranchSelected:
		if ni, ok := action.Branch.Contents.(models.NodeContentsInterface); ok {
			s.selected = ni.GetNode()
		} else {
			s.selected = nil
		}
	case *actions.EditorFocus:
		payload.Notify(action.Editor.Node, NodeFocus)
	}
	return true
}

func mutateDeleteNode(ctx context.Context, n *node.Node, p *node.Node, b *node.Node) error {
	*b = *n.Backup()
	switch p.Type.NativeJsonType() {
	case json.J_MAP:
		if err := p.DeleteMapChild(n.Key); err != nil {
			return kerr.Wrap("BUUOWYSJNG", err)
		}
	case json.J_ARRAY:
		if err := p.DeleteArrayChild(n.Index); err != nil {
			return kerr.Wrap("RWFQSINACH", err)
		}
	case json.J_OBJECT:
		if err := p.DeleteObjectChild(ctx, n.Key); err != nil {
			return kerr.Wrap("XGVEXEOBUP", err)
		}
	}
	return nil
}

func mutateRestoreNode(ctx context.Context, n *node.Node, p *node.Node, b *node.Node) error {
	if err := n.Restore(ctx, b); err != nil {
		return kerr.Wrap("EVSGQSPUPT", err)
	}
	switch p.Type.NativeJsonType() {
	case json.J_MAP:
		// don't have to call n.InitialiseMapItem because the node is already
		// initialized
		if err := n.AddToMap(ctx, p, n.Key, true); err != nil {
			return kerr.Wrap("TOPLOONYCL", err)
		}
	case json.J_ARRAY:
		// don't have to call n.InitialiseArrayItem because the node is already
		// initialized
		if err := n.AddToArray(ctx, p, n.Index, true); err != nil {
			return kerr.Wrap("WFXSQYOEAY", err)
		}
	case json.J_OBJECT:
		// don't have to call n.InitialiseObjectField because the node is
		// already initialized
		if err := n.AddToObject(ctx, p, n.Rule, n.Key, true); err != nil {
			return kerr.Wrap("QMBJQMLOCY", err)
		}
	}
	return nil
}

func mutateAddNode(ctx context.Context, n *node.Node, p *node.Node, key string, index int, t *system.Type) error {
	switch p.Type.NativeJsonType() {
	case json.J_ARRAY:
		if err := n.InitialiseArrayItem(ctx, p, index); err != nil {
			return kerr.Wrap("QLBGMSQENC", err)
		}
		if err := n.AddToArray(ctx, p, index, true); err != nil {
			return kerr.Wrap("PLEJOTCSGH", err)
		}
	case json.J_MAP:
		if err := n.InitialiseMapItem(ctx, p, key); err != nil {
			return kerr.Wrap("KRTGPFYWIH", err)
		}
		if err := n.AddToMap(ctx, p, key, true); err != nil {
			return kerr.Wrap("UEPLLMTLDB", err)
		}
	}
	if err := n.SetValueZero(ctx, false, t); err != nil {
		return kerr.Wrap("NLSRNQGLLW", err)
	}
	return nil
}
