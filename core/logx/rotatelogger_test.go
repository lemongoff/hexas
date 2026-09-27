package logx

import (
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/lemongoff/hexas/core/fs"
	"github.com/lemongoff/hexas/core/stringx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDailyRotateRuleMarkRotated(t *testing.T) {
	t.Run("daily rule", func(t *testing.T) {
		var rule DailyRotateRule
		rule.MarkRotated()
		assert.Equal(t, getNowDate(), rule.rotatedTime)
	})

	t.Run("daily rule", func(t *testing.T) {
		rule := DefaultRotateRule("test", "-", 1, false)
		_, ok := rule.(*DailyRotateRule)
		assert.True(t, ok)
	})
}

func TestDailyRotateRuleOutdatedFiles(t *testing.T) {
	t.Run("no files", func(t *testing.T) {
		var rule DailyRotateRule
		assert.Empty(t, rule.OutdatedFiles())
		rule.days = 1
		assert.Empty(t, rule.OutdatedFiles())
		rule.gzip = true
		assert.Empty(t, rule.OutdatedFiles())
	})

	t.Run("bad files", func(t *testing.T) {
		rule := DailyRotateRule{
			filename: "[a-z",
		}
		assert.Empty(t, rule.OutdatedFiles())
		rule.days = 1
		assert.Empty(t, rule.OutdatedFiles())
		rule.gzip = true
		assert.Empty(t, rule.OutdatedFiles())
	})

	t.Run("temp files", func(t *testing.T) {
		dir := t.TempDir()
		boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay) * 2).Format(time.DateOnly)
		f1, err := os.CreateTemp(dir, "go-zero-test-"+boundary)
		require.NoError(t, err)
		require.NoError(t, f1.Close())
		f2, err := os.CreateTemp(dir, "go-zero-test-"+boundary)
		require.NoError(t, err)
		require.NoError(t, f2.Close())
		rule := DailyRotateRule{
			filename: filepath.Join(dir, "go-zero-test-"),
			days:     1,
		}
		assert.ElementsMatch(t, []string{f1.Name(), f2.Name()}, rule.OutdatedFiles())
		// Forward slashes are accepted on Windows too; Glob returns native paths.
		rule.filename = filepath.ToSlash(rule.filename)
		assert.ElementsMatch(t, []string{f1.Name(), f2.Name()}, rule.OutdatedFiles())
	})
}

func TestDailyRotateRuleShallRotate(t *testing.T) {
	var rule DailyRotateRule
	rule.rotatedTime = time.Now().Add(time.Hour * 24).Format(time.DateOnly)
	assert.True(t, rule.ShallRotate(0))
}

func TestSizeLimitRotateRuleMarkRotated(t *testing.T) {
	t.Run("size limit rule", func(t *testing.T) {
		var rule SizeLimitRotateRule
		rule.MarkRotated()
		assert.Equal(t, getNowDateInRFC3339Format(), rule.rotatedTime)
	})

	t.Run("size limit rule", func(t *testing.T) {
		rule := NewSizeLimitRotateRule("foo", "-", 1, 1, 1, false)
		rule.MarkRotated()
		assert.Equal(t, getNowDateInRFC3339Format(), rule.(*SizeLimitRotateRule).rotatedTime)
	})
}

func TestSizeLimitRotateRuleOutdatedFiles(t *testing.T) {
	t.Run("no files", func(t *testing.T) {
		var rule SizeLimitRotateRule
		assert.Empty(t, rule.OutdatedFiles())
		rule.days = 1
		assert.Empty(t, rule.OutdatedFiles())
		rule.gzip = true
		assert.Empty(t, rule.OutdatedFiles())
		rule.maxBackups = 0
		assert.Empty(t, rule.OutdatedFiles())
	})

	t.Run("bad files", func(t *testing.T) {
		rule := SizeLimitRotateRule{
			DailyRotateRule: DailyRotateRule{
				filename: "[a-z",
			},
		}
		assert.Empty(t, rule.OutdatedFiles())
		rule.days = 1
		assert.Empty(t, rule.OutdatedFiles())
		rule.gzip = true
		assert.Empty(t, rule.OutdatedFiles())
	})

	t.Run("temp files", func(t *testing.T) {
		dir := t.TempDir()
		boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay) * 2).Format(time.DateOnly)
		f1, err := os.CreateTemp(dir, "go-zero-test-"+boundary)
		require.NoError(t, err)
		require.NoError(t, f1.Close())
		f2, err := os.CreateTemp(dir, "go-zero-test-"+boundary)
		require.NoError(t, err)
		require.NoError(t, f2.Close())
		boundary1 := time.Now().Add(time.Hour * time.Duration(hoursPerDay) * 2).Format(time.DateOnly)
		f3, err := os.CreateTemp(dir, "go-zero-test-"+boundary1)
		require.NoError(t, err)
		require.NoError(t, f3.Close())
		rule := SizeLimitRotateRule{
			DailyRotateRule: DailyRotateRule{
				filename: filepath.Join(dir, "go-zero-test-"),
				days:     1,
			},
			maxBackups: 3,
		}
		assert.NotEmpty(t, rule.OutdatedFiles())
	})

	t.Run("no backups", func(t *testing.T) {
		dir := t.TempDir()
		boundary := time.Now().Add(-time.Hour * time.Duration(hoursPerDay) * 2).Format(time.DateOnly)
		f1, err := os.CreateTemp(dir, "go-zero-test-"+boundary)
		require.NoError(t, err)
		require.NoError(t, f1.Close())
		f2, err := os.CreateTemp(dir, "go-zero-test-"+boundary)
		require.NoError(t, err)
		require.NoError(t, f2.Close())
		boundary1 := time.Now().Add(time.Hour * time.Duration(hoursPerDay) * 2).Format(time.DateOnly)
		f3, err := os.CreateTemp(dir, "go-zero-test-"+boundary1)
		require.NoError(t, err)
		require.NoError(t, f3.Close())
		rule := SizeLimitRotateRule{
			DailyRotateRule: DailyRotateRule{
				filename: filepath.Join(dir, "go-zero-test-"),
				days:     1,
			},
		}
		assert.NotEmpty(t, rule.OutdatedFiles())

		logger := new(RotateLogger)
		logger.rule = &rule
		logger.maybeDeleteOutdatedFiles()
		assert.Empty(t, rule.OutdatedFiles())
	})
}

func TestSizeLimitRotateRuleShallRotate(t *testing.T) {
	var rule SizeLimitRotateRule
	rule.rotatedTime = time.Now().Add(time.Hour * 24).Format(fileTimeFormat)
	rule.maxSize = 0
	assert.False(t, rule.ShallRotate(0))
	rule.maxSize = 100
	assert.False(t, rule.ShallRotate(0))
	assert.True(t, rule.ShallRotate(101*megaBytes))
}

func TestRotateLoggerClose(t *testing.T) {
	t.Run("close", func(t *testing.T) {
		filename, err := fs.TempFilenameWithText("foo")
		assert.Nil(t, err)
		if len(filename) > 0 {
			defer os.Remove(filename)
		}
		logger, err := NewLogger(filename, new(DailyRotateRule), false)
		assert.Nil(t, err)
		_, err = logger.Write([]byte("foo"))
		assert.Nil(t, err)
		assert.Nil(t, logger.Close())
	})

	t.Run("close and write", func(t *testing.T) {
		logger := new(RotateLogger)
		logger.done = make(chan struct{})
		close(logger.done)
		_, err := logger.Write([]byte("foo"))
		assert.ErrorIs(t, err, ErrLogFileClosed)
	})

	t.Run("close without losing logs", func(t *testing.T) {
		text := "foo"
		filename, err := fs.TempFilenameWithText(text)
		assert.Nil(t, err)
		if len(filename) > 0 {
			defer os.Remove(filename)
		}
		logger, err := NewLogger(filename, new(DailyRotateRule), false)
		assert.Nil(t, err)
		msg := []byte("foo")
		n := 100
		for i := 0; i < n; i++ {
			_, err = logger.Write(msg)
			assert.Nil(t, err)
		}
		assert.Nil(t, logger.Close())
		bs, err := os.ReadFile(filename)
		assert.Nil(t, err)
		assert.Equal(t, len(msg)*n+len(text), len(bs))
	})
}

func TestRotateLoggerGetBackupFilename(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}
	logger, err := NewLogger(filename, new(DailyRotateRule), false)
	assert.Nil(t, err)
	assert.True(t, len(logger.getBackupFilename()) > 0)
	logger.backup = ""
	assert.True(t, len(logger.getBackupFilename()) > 0)
}

func TestRotateLoggerMayCompressFile(t *testing.T) {
	old := os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)
	defer func() {
		os.Stdout = old
	}()

	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}
	logger, err := NewLogger(filename, new(DailyRotateRule), false)
	assert.Nil(t, err)
	logger.maybeCompressFile(filename)
	_, err = os.Stat(filename)
	assert.Nil(t, err)
}

func TestRotateLoggerMayCompressFileTrue(t *testing.T) {
	logger := newTestRotateLogger(t, false)
	require.NoError(t, logger.Close())
	logger.maybeCompressFile(logger.filename)
	require.NoFileExists(t, logger.filename)
	assertCompressedLog(t, logger.filename+gzipExt, "foo")
}

func TestRotateLoggerRotate(t *testing.T) {
	testRotateLoggerRotation(t, false)
}

func TestRotateLoggerWrite(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	rule := new(DailyRotateRule)
	logger, err := NewLogger(filename, rule, true)
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer func() {
			os.Remove(logger.getBackupFilename())
			os.Remove(filepath.Base(logger.getBackupFilename()) + ".gz")
		}()
	}
	// the following write calls cannot be changed to Write, because of DATA RACE.
	logger.write([]byte(`foo`))
	rule.rotatedTime = time.Now().Add(-time.Hour * 24).Format(time.DateOnly)
	logger.write([]byte(`bar`))
	logger.Close()
	logger.write([]byte(`baz`))
}

func TestLogWriterClose(t *testing.T) {
	assert.Nil(t, newLogWriter(nil).Close())
}

func TestRotateLoggerWithSizeLimitRotateRuleClose(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}
	logger, err := NewLogger(filename, new(SizeLimitRotateRule), false)
	assert.Nil(t, err)
	_ = logger.Close()
}

func TestRotateLoggerGetBackupWithSizeLimitRotateRuleFilename(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}
	logger, err := NewLogger(filename, new(SizeLimitRotateRule), false)
	assert.Nil(t, err)
	assert.True(t, len(logger.getBackupFilename()) > 0)
	logger.backup = ""
	assert.True(t, len(logger.getBackupFilename()) > 0)
}

func TestRotateLoggerWithSizeLimitRotateRuleMayCompressFile(t *testing.T) {
	old := os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)
	defer func() {
		os.Stdout = old
	}()

	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}
	logger, err := NewLogger(filename, new(SizeLimitRotateRule), false)
	assert.Nil(t, err)
	logger.maybeCompressFile(filename)
	_, err = os.Stat(filename)
	assert.Nil(t, err)
}

func TestRotateLoggerWithSizeLimitRotateRuleMayCompressFileTrue(t *testing.T) {
	logger := newTestRotateLogger(t, true)
	require.NoError(t, logger.Close())
	logger.maybeCompressFile(logger.filename)
	require.NoFileExists(t, logger.filename)
	assertCompressedLog(t, logger.filename+gzipExt, "foo")
}

func TestRotateLoggerWithSizeLimitRotateRuleMayCompressFileFailed(t *testing.T) {
	old := os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)
	defer func() {
		os.Stdout = old
	}()

	filename := stringx.RandId()
	logger, err := NewLogger(filename, new(SizeLimitRotateRule), true)
	defer os.Remove(filename)
	if assert.NoError(t, err) {
		assert.NotPanics(t, func() {
			logger.maybeCompressFile(stringx.RandId())
		})
	}
}

func TestRotateLoggerWithSizeLimitRotateRuleRotate(t *testing.T) {
	testRotateLoggerRotation(t, true)
}

func newTestRotateLogger(t *testing.T, sizeLimit bool) *RotateLogger {
	t.Helper()
	previousFormat := fileTimeFormat
	fileTimeFormat = defaultFileTimeFormat()
	t.Cleanup(func() { fileTimeFormat = previousFormat })
	filename := filepath.Join(t.TempDir(), "nested", "test.log")
	rule := DefaultRotateRule(filename, "-", 0, true)
	if sizeLimit {
		rule = NewSizeLimitRotateRule(filename, "-", 0, 1, 0, true)
	}
	logger, err := NewLogger(filename, rule, true)
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, logger.Close()) })
	logger.write([]byte("foo"))
	return logger
}

func testRotateLoggerRotation(t *testing.T, sizeLimit bool) {
	t.Helper()
	logger := newTestRotateLogger(t, sizeLimit)
	backup := logger.getBackupFilename()
	require.Equal(t, filepath.Dir(logger.filename), filepath.Dir(backup))
	require.NoError(t, logger.rotate())
	// Wait for compression to close its files before checking contents or cleaning up.
	require.Eventually(t, func() bool {
		_, err := os.Stat(backup)
		return errors.Is(err, os.ErrNotExist)
	}, 5*time.Second, time.Millisecond)
	assertCompressedLog(t, backup+gzipExt, "foo")
	logger.write([]byte("bar"))
	require.NoError(t, logger.Close())
	data, err := os.ReadFile(logger.filename)
	require.NoError(t, err)
	assert.Equal(t, "bar", string(data))
}

func assertCompressedLog(t *testing.T, filename, want string) {
	t.Helper()
	f, err := os.Open(filename)
	require.NoError(t, err)
	defer f.Close()
	r, err := gzip.NewReader(f)
	require.NoError(t, err)
	defer r.Close()
	data, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.Equal(t, want, string(data))
}

func TestRotateLoggerWithSizeLimitRotateRuleWrite(t *testing.T) {
	filename, err := fs.TempFilenameWithText("foo")
	assert.Nil(t, err)
	rule := new(SizeLimitRotateRule)
	logger, err := NewLogger(filename, rule, true)
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer func() {
			os.Remove(logger.getBackupFilename())
			os.Remove(filepath.Base(logger.getBackupFilename()) + ".gz")
		}()
	}
	// the following write calls cannot be changed to Write, because of DATA RACE.
	logger.write([]byte(`foo`))
	rule.rotatedTime = time.Now().Add(-time.Hour * 24).Format(time.DateOnly)
	logger.write([]byte(`bar`))
	logger.Close()
	logger.write([]byte(`baz`))
}

func TestGzipFile(t *testing.T) {
	err := errors.New("any error")

	t.Run("gzip file open failed", func(t *testing.T) {
		fsys := &fakeFileSystem{
			openFn: func(name string) (*os.File, error) {
				return nil, err
			},
		}
		assert.ErrorIs(t, err, gzipFile("any", fsys))
		assert.False(t, fsys.Removed())
	})

	t.Run("gzip file create failed", func(t *testing.T) {
		fsys := &fakeFileSystem{
			createFn: func(name string) (*os.File, error) {
				return nil, err
			},
		}
		assert.ErrorIs(t, err, gzipFile("any", fsys))
		assert.False(t, fsys.Removed())
	})

	t.Run("gzip file copy failed", func(t *testing.T) {
		fsys := &fakeFileSystem{
			copyFn: func(writer io.Writer, reader io.Reader) (int64, error) {
				return 0, err
			},
		}
		assert.ErrorIs(t, err, gzipFile("any", fsys))
		assert.False(t, fsys.Removed())
	})

	t.Run("gzip file last close failed", func(t *testing.T) {
		var called int32
		fsys := &fakeFileSystem{
			closeFn: func(closer io.Closer) error {
				if atomic.AddInt32(&called, 1) > 2 {
					return err
				}
				return nil
			},
		}
		assert.NoError(t, gzipFile("any", fsys))
		assert.True(t, fsys.Removed())
	})

	t.Run("gzip file remove failed", func(t *testing.T) {
		fsys := &fakeFileSystem{
			removeFn: func(name string) error {
				return err
			},
		}
		assert.Error(t, err, gzipFile("any", fsys))
		assert.True(t, fsys.Removed())
	})

	t.Run("gzip file everything ok", func(t *testing.T) {
		fsys := &fakeFileSystem{}
		assert.NoError(t, gzipFile("any", fsys))
		assert.True(t, fsys.Removed())
	})
}

func TestRotateLogger_WithExistingFile(t *testing.T) {
	const body = "foo"
	filename, err := fs.TempFilenameWithText(body)
	assert.Nil(t, err)
	if len(filename) > 0 {
		defer os.Remove(filename)
	}

	rule := NewSizeLimitRotateRule(filename, "-", 1, 100, 3, false)
	logger, err := NewLogger(filename, rule, false)
	assert.Nil(t, err)
	assert.Equal(t, int64(len(body)), logger.currentSize)
	assert.Nil(t, logger.Close())
}

func BenchmarkRotateLogger(b *testing.B) {
	filename := "./test.log"
	filename2 := "./test2.log"
	dailyRotateRuleLogger, err1 := NewLogger(
		filename,
		DefaultRotateRule(
			filename,
			backupFileDelimiter,
			1,
			true,
		),
		true,
	)
	if err1 != nil {
		b.Logf("Failed to new daily rotate rule logger: %v", err1)
		b.FailNow()
	}
	sizeLimitRotateRuleLogger, err2 := NewLogger(
		filename2,
		NewSizeLimitRotateRule(
			filename,
			backupFileDelimiter,
			1,
			100,
			10,
			true,
		),
		true,
	)
	if err2 != nil {
		b.Logf("Failed to new size limit rotate rule logger: %v", err1)
		b.FailNow()
	}
	defer func() {
		dailyRotateRuleLogger.Close()
		sizeLimitRotateRuleLogger.Close()
		os.Remove(filename)
		os.Remove(filename2)
	}()

	b.Run("daily rotate rule", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dailyRotateRuleLogger.write([]byte("testing\ntesting\n"))
		}
	})
	b.Run("size limit rotate rule", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sizeLimitRotateRuleLogger.write([]byte("testing\ntesting\n"))
		}
	})
}

type fakeFileSystem struct {
	removed  int32
	closeFn  func(closer io.Closer) error
	copyFn   func(writer io.Writer, reader io.Reader) (int64, error)
	createFn func(name string) (*os.File, error)
	openFn   func(name string) (*os.File, error)
	removeFn func(name string) error
}

func (f *fakeFileSystem) Close(closer io.Closer) error {
	if f.closeFn != nil {
		return f.closeFn(closer)
	}
	return nil
}

func (f *fakeFileSystem) Copy(writer io.Writer, reader io.Reader) (int64, error) {
	if f.copyFn != nil {
		return f.copyFn(writer, reader)
	}
	return 0, nil
}

func (f *fakeFileSystem) Create(name string) (*os.File, error) {
	if f.createFn != nil {
		return f.createFn(name)
	}
	return nil, nil
}

func (f *fakeFileSystem) Open(name string) (*os.File, error) {
	if f.openFn != nil {
		return f.openFn(name)
	}
	return nil, nil
}

func (f *fakeFileSystem) Remove(name string) error {
	atomic.AddInt32(&f.removed, 1)

	if f.removeFn != nil {
		return f.removeFn(name)
	}
	return nil
}

func (f *fakeFileSystem) Removed() bool {
	return atomic.LoadInt32(&f.removed) > 0
}
