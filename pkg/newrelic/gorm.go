package newrelic

import (
	"context"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormNewRelicLogger wraps GORM logger with New Relic monitoring.
type GormNewRelicLogger struct {
	logger.Interface
	app *newrelic.Application
}

// NewGormLogger creates a new GORM logger with New Relic integration.
func NewGormLogger(baseLogger logger.Interface, app *newrelic.Application) *GormNewRelicLogger {
	return &GormNewRelicLogger{
		Interface: baseLogger,
		app:       app,
	}
}

// Trace adds New Relic database segment to GORM queries.
func (l *GormNewRelicLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.app != nil {
		// Extract New Relic transaction from context
		if txn := newrelic.FromContext(ctx); txn != nil {
			sql, _ := fc()
			segment := newrelic.DatastoreSegment{
				Product:    newrelic.DatastorePostgres, // Change to DatastoreSQLite for SQLite
				Collection: "query",
				Operation:  "exec",
			}
			if sql != "" {
				segment.ParameterizedQuery = sql
			}
			segment.StartTime = txn.StartSegmentNow()
			segment.End()
		}
	}

	// Call the original trace method
	l.Interface.Trace(ctx, begin, fc, err)
}

// AddNewRelicToGorm adds New Relic callbacks to GORM instance.
func AddNewRelicToGorm(db *gorm.DB, app *newrelic.Application) error {
	if app == nil {
		return nil
	}

	// Add New Relic callback for create operations
	if err := db.Callback().Create().Before("gorm:create").Register("newrelic:before_create", beforeCreate(app)); err != nil {
		return err
	}
	if err := db.Callback().Create().After("gorm:create").Register("newrelic:after_create", afterCreate); err != nil {
		return err
	}

	// Add New Relic callback for query operations
	if err := db.Callback().Query().Before("gorm:query").Register("newrelic:before_query", beforeQuery(app)); err != nil {
		return err
	}
	if err := db.Callback().Query().After("gorm:query").Register("newrelic:after_query", afterQuery); err != nil {
		return err
	}

	// Add New Relic callback for update operations
	if err := db.Callback().Update().Before("gorm:update").Register("newrelic:before_update", beforeUpdate(app)); err != nil {
		return err
	}
	if err := db.Callback().Update().After("gorm:update").Register("newrelic:after_update", afterUpdate); err != nil {
		return err
	}

	// Add New Relic callback for delete operations
	if err := db.Callback().Delete().Before("gorm:delete").Register("newrelic:before_delete", beforeDelete(app)); err != nil {
		return err
	}
	if err := db.Callback().Delete().After("gorm:delete").Register("newrelic:after_delete", afterDelete); err != nil {
		return err
	}

	return nil
}

func beforeCreate(app *newrelic.Application) func(*gorm.DB) {
	return func(db *gorm.DB) {
		startDatastoreSegment(db, "CREATE")
	}
}

func afterCreate(db *gorm.DB) {
	endDatastoreSegment(db)
}

func beforeQuery(app *newrelic.Application) func(*gorm.DB) {
	return func(db *gorm.DB) {
		startDatastoreSegment(db, "SELECT")
	}
}

func afterQuery(db *gorm.DB) {
	endDatastoreSegment(db)
}

func beforeUpdate(app *newrelic.Application) func(*gorm.DB) {
	return func(db *gorm.DB) {
		startDatastoreSegment(db, "UPDATE")
	}
}

func afterUpdate(db *gorm.DB) {
	endDatastoreSegment(db)
}

func beforeDelete(app *newrelic.Application) func(*gorm.DB) {
	return func(db *gorm.DB) {
		startDatastoreSegment(db, "DELETE")
	}
}

func afterDelete(db *gorm.DB) {
	endDatastoreSegment(db)
}

func startDatastoreSegment(db *gorm.DB, operation string) {
	if txn := newrelic.FromContext(db.Statement.Context); txn != nil {
		segment := &newrelic.DatastoreSegment{
			Product:   newrelic.DatastorePostgres, // Change based on your database
			Operation: operation,
		}
		segment.StartTime = txn.StartSegmentNow()
		db.Set("newrelic:segment", segment)
	}
}

func endDatastoreSegment(db *gorm.DB) {
	if segment, exists := db.Get("newrelic:segment"); exists {
		if datastoreSegment, ok := segment.(*newrelic.DatastoreSegment); ok {
			if db.Statement.Table != "" {
				datastoreSegment.Collection = db.Statement.Table
			}
			datastoreSegment.End()
		}
	}
}
