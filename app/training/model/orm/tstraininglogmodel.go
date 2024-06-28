package orm

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TsTrainingLogModel = (*customTsTrainingLogModel)(nil)

type (
	// TsTrainingLogModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTsTrainingLogModel.
	TsTrainingLogModel interface {
		tsTrainingLogModel
		withSession(session sqlx.Session) TsTrainingLogModel
	}

	customTsTrainingLogModel struct {
		*defaultTsTrainingLogModel
	}
)

// NewTsTrainingLogModel returns a model for the database table.
func NewTsTrainingLogModel(conn sqlx.SqlConn) TsTrainingLogModel {
	return &customTsTrainingLogModel{
		defaultTsTrainingLogModel: newTsTrainingLogModel(conn),
	}
}

func (m *customTsTrainingLogModel) withSession(session sqlx.Session) TsTrainingLogModel {
	return NewTsTrainingLogModel(sqlx.NewSqlConnFromSession(session))
}
