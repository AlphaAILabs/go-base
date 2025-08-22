package postgres

import (
	"testing"
)

func TestNew(t *testing.T) {
	config := &Config{
		User: "alpha_labs_test",
		Pass: "alpha_dsdh2377293cy2c0238ff23ch",
		Host: "127.0.0.1",
		Port: 5432,
		Db:   "alpahlabs_product_test",
	}

	t.Run("testpostgres", func(t *testing.T) {
		pg, err := New(config)
		t.Logf("err:%s", err)
		t.Log("config: ", pg.config)
		count := -1
		err = pg.Client.Table("test").Count(&count).Error
		t.Logf("count: %d", count)
		t.Logf("count err:%s", err)
	})
}
