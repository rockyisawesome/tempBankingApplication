package repositories

import (
	"context"
	"fmt"
	"ledgerservice/database"
	"ledgerservice/models"
	"time"

	"github.com/hashicorp/go-hclog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TransactionRepo struct {
	mgdb  database.Database
	loggs *hclog.Logger
}

func NewTransactionRepository(mgdb database.Database, lobbs *hclog.Logger) Repository {
	return &TransactionRepo{
		loggs: lobbs,
		mgdb:  mgdb,
	}
}

func (t *TransactionRepo) InsertTransaction(ctx context.Context, ledger models.TransactionLedger) error {

	// update the value of the document
	ledger.ID = bson.NewObjectID() // Set a new ObjectID
	ledger.CreatedAt = time.Now()  // Set creation timestamp
	ledger.Status = "completed"

	// Insert into ledger
	obid, err := t.mgdb.InsertTransaction(ctx, ledger)
	if err != nil {
		(*t.loggs).Error("Error inserting transaction", "Error", err)
		return err
	}

	fmt.Println("Transaction inserted successfully with ID: ", obid)
	return nil
}
