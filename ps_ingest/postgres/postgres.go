package postgres

import (
	"database/sql"
	"fmt"
	datatypes "ps_ingest/dataTypes"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type PgClient interface {
	Connect() error
	InsertGame(game datatypes.Game, descEmbedding [][]float64) error
	Close()
}

type client struct {
	user     string
	password string
	host     string
	database string
	db       *sql.DB
}

func (c *client) Connect() error {
	log.Info().Str("database", c.database).Msg("Connecting to postgres db")
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=disable", c.user, c.password, c.host, c.database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to postgres")
		return err
	}

	c.db = db
	return nil
}

func (c *client) InsertGame(game datatypes.Game, descEmbedding [][]float64) error {
	_, err := c.db.Exec(`
		INSERT INTO games (name, description, price, original_price, url, rating, rating_sum, expiration, description_embedding)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`,
		game.Name,
		game.Description,
		game.Price,
		game.OriginalPrice,
		game.URL,
		game.Rating,
		game.RatingsSum,
		game.Expiration,
		// TODO: remodel db to store whole description embedding
		c.toVectorInput(descEmbedding[0]),
	)

	return err
}

func (c *client) toVectorInput(emb []float64) string {
	stringValues := []string{}
	for _, e := range emb {
		stringValues = append(stringValues, strconv.FormatFloat(e, 'E', -1, 64))
	}

	return fmt.Sprintf("[%s]", strings.Join(stringValues, ","))
}

func (c *client) Close() {
	log.Info().Msg("Closing postgres connection")
	if c.db == nil {
		return
	}

	err := c.db.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close postgres connection")
	}
}

type PgConfig struct {
	User     string
	Password string
	Host     string
	Database string
}

func CreatePgClient(cfg PgConfig) PgClient {
	return &client{cfg.User, cfg.Password, cfg.Host, cfg.Database, nil}
}
