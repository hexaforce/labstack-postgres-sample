package crud

import "github.com/go-pg/pg/v10"

func CreateComment(db *pg.DB, req *Comment) (*Comment, error) {
	_, err := db.Model(req).Insert()
	if err != nil {
		return nil, err
	}

	comment := &Comment{}

	err = db.Model(comment).
		Relation("User").
		Where("comment.id = ?", req.ID).
		Select()

	return comment, err
}

func GetComment(db *pg.DB, commentID string) (*Comment, error) {
	comment := &Comment{}

	err := db.Model(comment).
		Relation("User").
		Where("comment.id = ?", commentID).
		Select()

	return comment, err
}

func GetComments(db *pg.DB) ([]*Comment, error) {
	comments := make([]*Comment, 0)

	err := db.Model(&comments).
		Relation("User").
		Select()

	return comments, err
}

func UpdateComment(db *pg.DB, req *Comment) (*Comment, error) {
	_, err := db.Model(req).
		WherePK().
		Update()
	if err != nil {
		return nil, err
	}

	comment := &Comment{}

	err = db.Model(comment).
		Relation("User").
		Where("comment.id = ?", req.ID).
		Select()

	return comment, err
}

func DeleteComment(db *pg.DB, commentID int64) error {
	comment := &Comment{}

	err := db.Model(comment).
		Relation("User").
		Where("comment.id = ?", commentID).
		Select()
	if err != nil {
		return err
	}

	_, err = db.Model(comment).WherePK().Delete()

	return err
}
