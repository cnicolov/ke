package stores

import (
	"golang.org/x/net/context"
	"kego.io/editor/client/actions"
	"kego.io/editor/client/models"
	"kego.io/flux"
	"kego.io/system/node"
)

type EditorStore struct {
	*flux.Store
	ctx context.Context
	app *App

	editors map[*node.Node]*models.EditorModel
}

type editorNotif string

func (b editorNotif) IsNotif() {}

const (
	EditorChanged            editorNotif = "EditorChanged"
	EditorSelected           editorNotif = "EditorSelected"
	EditorLoaded             editorNotif = "EditorLoaded"
	EditorAdded              editorNotif = "EditorAdded"
	EditorInitialStateLoaded editorNotif = "EditorInitialStateLoaded"
	EditorArrayOrderChanged  editorNotif = "EditorArrayOrderChanged"
)

func NewEditorStore(ctx context.Context) *EditorStore {
	s := &EditorStore{
		Store:   &flux.Store{},
		ctx:     ctx,
		app:     FromContext(ctx),
		editors: map[*node.Node]*models.EditorModel{},
	}
	s.Init(s)
	return s
}

func (s *EditorStore) Get(node *node.Node) *models.EditorModel {
	if node == nil {
		return nil
	}
	return s.editors[node]
}

func (s *EditorStore) Handle(payload *flux.Payload) bool {
	switch action := payload.Action.(type) {
	case *actions.InitialState:
		payload.Wait(s.app.Package, s.app.Types, s.app.Data)
		s.AddEditorsRecursively(s.app.Package.Node())
		for _, n := range s.app.Types.All() {
			s.AddEditorsRecursively(n)
		}
		s.Notify(nil, EditorInitialStateLoaded)
	case *actions.LoadSourceSuccess:
		ni, ok := action.Branch.Contents.(models.NodeContentsInterface)
		if !ok {
			break
		}
		n := ni.GetNode()
		e := s.AddEditorsRecursively(n)
		s.Notify(e, EditorLoaded)
	case *actions.InitializeNode:
		payload.Wait(s.app.Branches)
		e := s.AddEditorsRecursively(action.Node)
		s.Notify(e, EditorAdded)
	case *actions.BranchSelected:
		payload.Wait(s.app.Nodes)
		if e := s.Get(s.app.Nodes.Selected()); e != nil {
			s.Notify(e, EditorSelected)
		}
		s.Notify(nil, EditorSelected)
	case *actions.DeleteNode:
		payload.Wait(s.app.Nodes)
		e, ok := s.editors[action.Node]
		if !ok {
			break
		}
		if action.Node.Parent.Type.IsNativeCollection() {
			delete(s.editors, action.Node)
		}
		s.Notify(e, EditorChanged)
	case *actions.ArrayOrder:
		payload.Wait(s.app.Branches)
		s.Notify(nil, EditorArrayOrderChanged)
	case *actions.EditorValueChange:
		action.Editor.TemporaryValue = action.Value
	}
	return true
}

func (s *EditorStore) AddEditorsRecursively(n *node.Node) *models.EditorModel {
	e := models.NewEditor(n)
	s.editors[n] = e
	for _, c := range n.Map {
		s.AddEditorsRecursively(c)
	}
	for _, c := range n.Array {
		s.AddEditorsRecursively(c)
	}
	return e
}
