package migrate

import (
	"gohub/pkg/console"
	"gohub/pkg/database"
	"gohub/pkg/file"
	"gorm.io/gorm"
	"io/ioutil"
)

type Migrator struct {
	Folder   string
	DB       *gorm.DB
	Migrator gorm.Migrator
}

type Migration struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;"`
	Migration string `gorm:"type:varchar(255);not null;unique;"`
	Batch     int
}

func NewMigrator() *Migrator {
	migrator := &Migrator{
		Folder:   "database/migrations",
		DB:       database.DB,
		Migrator: database.DB.Migrator(),
	}

	migrator.createMigrationsTable()

	return migrator
}

func (migrator *Migrator) createMigrationsTable() {
	migration := Migration{}

	if !migrator.Migrator.HasTable(migration) {
		migrator.Migrator.CreateTable(migration)
	}
}

func (migrator *Migrator) Up() {
	// 读取所有迁移文件, 确保按照时间排序
	migrateFiles := migrator.readAllMigrationFiles()

	// 获取当前批次的值
	batch := migrator.getBatch()

	// 获取所有迁移数据
	migrations := []Migration{}
	migrator.DB.Find(&migrations)

	// 通过此值判断数据库是否是最新
	runed := false

	// 对迁移文件遍历, 如果没有执行过, 就执行up 回调
	for _, mfile := range migrateFiles {
		if mfile.isNotMigrated(migrations) {
			migrator.runUpMigration(mfile, batch)
			runed = true
		}
	}

	if !runed {
		console.Success("database is up to date.")
	}
}

func (migrator *Migrator) Rollback() {
	// 获取最后一批次的迁移数据
	lastMigration := Migration{}
	migrator.DB.Order("id DESC").First(&lastMigration)
	migrations := []Migration{}
	migrator.DB.Where("batch = ?", lastMigration.Batch).Order("id DESC").Find(&migrations)

	// 回滚最后一批次的迁移
	if !migrator.rollbackMigrations(migrations) {
		console.Success("[migrations] table is empty, nothing to rollback.")
	}
}

func (migrator *Migrator) rollbackMigrations(migrations []Migration) bool {
	// 标记是否真的有执行了迁移回退的操作
	runed := false

	for _, _migration := range migrations {
		console.Warning("rollback " + _migration.Migration)

		// 执行迁移文件的 down 方法
		mfile := getMigrationFile(_migration.Migration)
		if mfile.Down != nil {
			mfile.Down(database.DB.Migrator(), database.SQLDB)
		}

		runed = true

		migrator.DB.Delete(&_migration)

		console.Success("finish " + mfile.FileName)
	}

	return runed
}

func (migrator *Migrator) getBatch() int {
	batch := 1

	lastMigration := Migration{}
	migrator.DB.Order("id DESC").First(&lastMigration)

	if lastMigration.ID > 0 {
		batch = lastMigration.Batch + 1
	}

	return batch
}

func (migrator *Migrator) readAllMigrationFiles() []MigrationFile {
	files, err := ioutil.ReadDir(migrator.Folder)
	console.ExitIf(err)

	var migrateFiles []MigrationFile

	for _, f := range files {
		// 去掉文件后缀 go
		fileName := file.FileNameWithoutExtension(f.Name())

		mfile := getMigrationFile(fileName)

		if len(mfile.FileName) > 0 {
			migrateFiles = append(migrateFiles, mfile)
		}
	}

	return migrationFiles
}

func (migrator *Migrator) runUpMigration(mfile MigrationFile, batch int) {
	// 执行 up 区块的 sql
	if mfile.Up != nil {
		// 友好提示
		console.Warning("migrating " + mfile.FileName)
		// 执行 up 方法
		mfile.Up(database.DB.Migrator(), database.SQLDB)
		// 提示已经迁移的那个文件
		console.Success("migrated " + mfile.FileName)
	}

	err := migrator.DB.Create(&Migration{
		Migration: mfile.FileName,
		Batch:     batch,
	}).Error
	console.ExitIf(err)
}

func (migrator *Migrator) Reset() {
	migrations := []Migration{}
	migrator.DB.Order("id DESC").Find(&migrations)

	if !migrator.rollbackMigrations(migrations) {
		console.Success("[migrations] table is empty, nothing to reset.")
	}
}

func (migrator *Migrator) Refresh() {
	migrator.Reset()
	migrator.Up()
}

func (migrator *Migrator) Fresh() {
	// 获取数据库名称, 用以提示
	dbname := database.CurrentDatabase()

	// 删除所有表
	err := database.DeleteAllTables()
	console.ExitIf(err)
	console.Success("clearup database " + dbname)

	// 重新创建 migrations 表
	migrator.createMigrationsTable()
	console.Success("['migrations] table created.")

	// 重新调用 up 命令
	migrator.Up()
}
