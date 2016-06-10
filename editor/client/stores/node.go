package stores

import (
	"github.com/davelondon/kerr"
	"golang.org/x/net/context"
	"kego.io/editor/client/actions"
	"kego.io/editor/client/models"
	"kego.io/flux"
	"kego.io/json"
	"kego.io/system/node"
)

type NodeStore struct {
	*flux.Store
	ctx context.Context
	app *App

	selected *node.Node

	addPop *models.AddPopModel
}

func (s *NodeStore) Selected() *node.Node {
	return s.selected
}

func (s *NodeStore) AddPop() *models.AddPopModel {
	return s.addPop
}

type nodeNotif string

func (b nodeNotif) IsNotif() {}

const (
	NodeInitialised nodeNotif = "NodeInitialised"
	AddPopChange    nodeNotif = "AddPopChange"
	NodeFocused     nodeNotif = "NodeFocused"
)

func NewNodeStore(ctx context.Context) *NodeStore {
	s := &NodeStore{
		Store: &flux.Store{},
		ctx:   ctx,
		app:   FromContext(ctx),
	}
	s.Init(s)
	return s
}

func (s *NodeStore) Handle(payload *flux.Payload) bool {
	switch action := payload.Action.(type) {
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
	case *actions.DeleteNode:
		if action.Node.Parent.Type.IsNativeCollection() {
			if action.Node.Index > -1 {
				action.Node.Parent.Array = append(action.Node.Parent.Array[:action.Node.Index], action.Node.Parent.Array[action.Node.Index+1:]...)
			} else {
				delete(action.Node.Parent.Map, action.Node.Key)
			}
		} else {
			action.Node.Type = nil
			action.Node.Array = []*node.Node{}
			action.Node.JsonType = json.J_NULL
			action.Node.Map = map[string]*node.Node{}
			action.Node.Missing = true
			action.Node.Null = true
			action.Node.Value = nil
			action.Node.ValueBool = false
			action.Node.ValueNumber = 0.0
			action.Node.ValueString = ""
		}
	case *actions.InitializeNode:
		if action.New {
			action.Node.Parent = action.Parent
			action.Node.Key = action.Key
			action.Node.Index = action.Index
			action.Node.Rule = action.Rule
			if action.Index > -1 {
				action.Parent.Array = append(action.Parent.Array, action.Node)
			} else {
				action.Parent.Map[action.Key] = action.Node
			}
		}
		if err := action.Node.InitialiseWithConcreteType(s.ctx, action.Type); err != nil {
			s.app.Fail <- kerr.Wrap("WWKUVDDLYU", err)
		}
	case *actions.OpenAddPop:

		s.addPop = &models.AddPopModel{
			Visible: true,
			Parent:  action.Parent,
			Node:    action.Node,
			Types:   action.Types,
		}
		s.Notify(nil, AddPopChange)
	case *actions.CloseAddPop:
		s.addPop = &models.AddPopModel{
			Visible: false,
		}
		s.Notify(nil, AddPopChange)
	case *actions.FocusNode:
		s.Notify(action.Node, NodeFocused)
	}
	return true
}
