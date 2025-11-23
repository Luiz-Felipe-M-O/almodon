package db

import(
	"fmt"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Essa função cria uma nova conexão com o banco de dados
func NewConnection(dataBaseUrl string) (*pgxpool.Pool, error){
	dbPoll, err := pgxpool.New(context.Background(), dataBaseUrl)

	if err != nil {
		dbPoll.Close()
		return nil, fmt.Errorf("Erro ao conectar ao Banco de dados: ", err)
	}

	return dbPoll, nil
}

func NewClinica(id int, nome string){
	
}