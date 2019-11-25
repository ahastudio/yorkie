package proxy

import (
	"github.com/hackerwins/yorkie/pkg/document/change"
	"github.com/hackerwins/yorkie/pkg/document/json/datatype"
	"github.com/hackerwins/yorkie/pkg/document/operation"
	"github.com/hackerwins/yorkie/pkg/document/time"
	"github.com/hackerwins/yorkie/pkg/log"
)

type TextProxy struct {
	*datatype.Text
	context *change.Context
}

func ProxyText(ctx *change.Context, text *datatype.Text) *TextProxy {
	rgaTreeSplit := datatype.NewRGATreeSplit()

	current := rgaTreeSplit.InitialHead()
	for _, textNode := range text.TextNodes() {
		current = rgaTreeSplit.InsertAfter(current, textNode.DeepCopy())
		insPrevID := textNode.InsPrevID()
		if insPrevID != nil {
			insPrevNode := rgaTreeSplit.FindTextNode(insPrevID)
			if insPrevNode == nil {
				log.Logger.Warn("insPrevNode should be presence")
			}
			current.SetInsPrev(insPrevNode)
		}
	}

	return NewTextProxy(ctx, rgaTreeSplit, text.CreatedAt())
}

func NewTextProxy(
	ctx *change.Context,
	rgaTreeSplit *datatype.RGATreeSplit,
	createdAt *time.Ticket,
) *TextProxy {
	return &TextProxy{
		Text:    datatype.NewText(rgaTreeSplit, createdAt),
		context: ctx,
	}
}

func (p *TextProxy) Edit(from, to int, content string) *TextProxy {
	if from > to {
		panic("from should be less than or equal to to")
	}
	fromPos, toPos := p.Text.FindBoundary(from, to)
	log.Logger.Debugf(
		"EDIT: f:%d->%s, t:%d->%s c:%s",
		from, fromPos.AnnotatedString(), to, toPos.AnnotatedString(), content,
	)

	ticket := p.context.IssueTimeTicket()
	_, maxCreationMapByActor := p.Text.Edit(
		fromPos,
		toPos,
		nil,
		content,
		ticket,
	)

	p.context.Push(operation.NewEdit(
		p.CreatedAt(),
		fromPos,
		toPos,
		maxCreationMapByActor,
		content,
		ticket,
	))

	return p
}