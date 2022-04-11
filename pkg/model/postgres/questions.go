package postgres

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"goquiz/pkg/model"
	"time"
)

func (m Model) InsertCategories(categs []model.Category) error {
	fmt.Println("Insert CategoriesModel")
	fmt.Printf("LEN: %d \n", len(categs))
	for _, categ := range categs {

		categID, err := m.CategoriesModel.InsertCategory(categ)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println("Successfully Write Category, ", categ, " with ID ", categID)
		fmt.Println("Total Question ", len(categ.Questions))

		// later refactor method - InsertQuestions
		for _, question := range categ.Questions {
			err := m.QuestionsModel.InsertQuestion(categID, question)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		}
		// TODO: add transaction rollback
		//if err != nil {
		//	// rollback transaction
		//	continue
		//}
		// commit transaction
	}
	return nil
}

func (m CategoriesModel) InsertCategory(cate model.Category) (int, error) {
	query := `INSERT INTO categories (name) values ($1) RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	categ := struct {
		id int
	}{}

	err := m.DB.QueryRowContext(ctx, query, cate.Name).Scan(&categ.id)
	fmt.Println("Category ID")
	fmt.Println(categ.id)
	if err != nil {
		return 0, err
	}
	return categ.id, nil
}

func (m QuestionsModel) InsertQuestion(categID int, q model.Question) error {
	args := []interface{}{categID, q.WebIndex, q.Text, pq.Array(q.AnsOptions), q.Codeblock, q.Answer.Option, q.Answer.Explanation, q.URL}
	// TODO: replace query with sql prepared statement
	query := `INSERT INTO questions (category_id,web_index,text,ans_options,code_block,correct_ans_opt,correct_ans_explanation,url) values ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}
	return nil
}
