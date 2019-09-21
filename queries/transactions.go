package queries

import (
	"github.com/DanielRenne/GoCore/core"
	"github.com/blauentag/blau/models/v1/model"
	"github.com/blauentag/blau/sessionFunctions"
	"github.com/pkg/errors"
)

type queryTransactions struct{}

func (self queryTransactions) Rollback(context session_functions.RequestContext, reason string, transactionId string) (err error) {
	var t model.Transaction
	err = model.Transactions.Query().ById(transactionId, &t)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	user, err := getUser(context)
	if err != nil {
		err = errors.Wrap(err, core.Debug.ErrLineAndFile(err))
		return
	}

	err = t.Rollback(user.Id.Hex(), reason)
	return
}
