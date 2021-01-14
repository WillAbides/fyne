package repository

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	"fyne.io/fyne/storage"
	"fyne.io/fyne/storage/repository"

	"github.com/stretchr/testify/assert"
)

func TestFileRepositoryRegistration(t *testing.T) {
	f := NewFileRepository("file")
	repository.Register("file", f)

	// this should never fail, and we assume it doesn't in other tests here
	// for brevity
	foo, err := storage.ParseURI("file:///foo")
	assert.Nil(t, err)

	// make sure we get the same repo back
	repo, err := repository.ForURI(foo)
	assert.Nil(t, err)
	assert.Equal(t, f, repo)
}

func TestFileRepositoryExists(t *testing.T) {
	dir, err := ioutil.TempDir("", "FyneInternalRepositoryFileTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	existsPath := path.Join(dir, "exists")
	notExistsPath := path.Join(dir, "notExists")

	err = ioutil.WriteFile(existsPath, []byte{1, 2, 3, 4}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	ex, err := storage.Exists(storage.NewFileURI(existsPath))
	assert.Nil(t, err)
	assert.True(t, ex)

	ex, err = storage.Exists(storage.NewFileURI(notExistsPath))
	assert.Nil(t, err)
	assert.False(t, ex)
}

func TestFileRepositoryReader(t *testing.T) {
	// Set up a temporary directory.
	dir, err := ioutil.TempDir("", "FyneInternalRepositoryFileTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create some files to test with.
	fooPath := path.Join(dir, "foo")
	barPath := path.Join(dir, "bar")
	bazPath := path.Join(dir, "baz")
	err = ioutil.WriteFile(fooPath, []byte{}, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(barPath, []byte{1, 2, 3}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Set up our repository - it's OK if we already registered it...
	f := NewFileRepository("file")
	repository.Register("file", f)

	// ...and some URIs - we know that they will not fail parsing
	foo := storage.NewFileURI(fooPath)
	bar := storage.NewFileURI(barPath)
	baz := storage.NewFileURI(bazPath)

	// Make sure we can read the empty file.
	fooReader, err := storage.Reader(foo)
	assert.Nil(t, err)
	fooData, err := ioutil.ReadAll(fooReader)
	assert.Equal(t, []byte{}, fooData)
	assert.Nil(t, err)

	// Make sure we can read the file with data.
	barReader, err := storage.Reader(bar)
	assert.Nil(t, err)
	barData, err := ioutil.ReadAll(barReader)
	assert.Equal(t, []byte{1, 2, 3}, barData)
	assert.Nil(t, err)

	// Make sure we get an error if the file doesn't exist.
	_, err = storage.Reader(baz)
	assert.NotNil(t, err)

	// Also test that CanRead returns the expected results.
	fooCanRead, err := storage.CanRead(foo)
	assert.True(t, fooCanRead)
	assert.Nil(t, err)

	barCanRead, err := storage.CanRead(bar)
	assert.True(t, barCanRead)
	assert.Nil(t, err)

	bazCanRead, err := storage.CanRead(baz)
	assert.False(t, bazCanRead)
	assert.Nil(t, err)
}

func TestFileRepositoryWriter(t *testing.T) {
	// Set up a temporary directory.
	dir, err := ioutil.TempDir("", "FyneInternalRepositoryFileTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create some files to test with.
	fooPath := path.Join(dir, "foo")
	barPath := path.Join(dir, "bar")
	bazPath := path.Join(dir, "baz")
	err = ioutil.WriteFile(fooPath, []byte{}, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(barPath, []byte{1, 2, 3}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Set up our repository - it's OK if we already registered it...
	f := NewFileRepository("file")
	repository.Register("file", f)

	// ...and some URIs - we know that they will not fail parsing
	foo := storage.NewFileURI(fooPath)
	bar := storage.NewFileURI(barPath)
	baz := storage.NewFileURI(bazPath)

	// write some data and assert there are no errors
	fooWriter, err := storage.Writer(foo)
	assert.Nil(t, err)
	assert.NotNil(t, fooWriter)

	barWriter, err := storage.Writer(bar)
	assert.Nil(t, err)
	assert.NotNil(t, barWriter)

	bazWriter, err := storage.Writer(baz)
	assert.Nil(t, err)
	assert.NotNil(t, bazWriter)

	n, err := fooWriter.Write([]byte{1, 2, 3, 4, 5})
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	n, err = barWriter.Write([]byte{6, 7, 8, 9})
	assert.Nil(t, err)
	assert.Equal(t, 4, n)

	n, err = bazWriter.Write([]byte{5, 4, 3, 2, 1, 0})
	assert.Nil(t, err)
	assert.Equal(t, 6, n)

	fooWriter.Close()
	barWriter.Close()
	bazWriter.Close()

	// now make sure we can read the data back correctly
	fooReader, err := storage.Reader(foo)
	assert.Nil(t, err)
	fooData, err := ioutil.ReadAll(fooReader)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, fooData)
	assert.Nil(t, err)

	barReader, err := storage.Reader(bar)
	assert.Nil(t, err)
	barData, err := ioutil.ReadAll(barReader)
	assert.Equal(t, []byte{6, 7, 8, 9}, barData)
	assert.Nil(t, err)

	bazReader, err := storage.Reader(baz)
	assert.Nil(t, err)
	bazData, err := ioutil.ReadAll(bazReader)
	assert.Equal(t, []byte{5, 4, 3, 2, 1, 0}, bazData)
	assert.Nil(t, err)

	// close the readers, since Windows won't let us delete things with
	// open handles to them
	fooReader.Close()
	barReader.Close()
	bazReader.Close()

	// now let's test deletion
	err = storage.Delete(foo)
	assert.Nil(t, err)

	err = storage.Delete(bar)
	assert.Nil(t, err)

	err = storage.Delete(baz)
	assert.Nil(t, err)

	fooExists, err := storage.Exists(foo)
	assert.False(t, fooExists)
	assert.Nil(t, err)

	barExists, err := storage.Exists(bar)
	assert.False(t, barExists)
	assert.Nil(t, err)

	bazExists, err := storage.Exists(baz)
	assert.False(t, bazExists)
	assert.Nil(t, err)
}

func TestFileRepositoryCanWrite(t *testing.T) {
	// Set up a temporary directory.
	dir, err := ioutil.TempDir("", "FyneInternalRepositoryFileTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create some files to test with.
	fooPath := path.Join(dir, "foo")
	barPath := path.Join(dir, "bar")
	bazPath := path.Join(dir, "baz")
	err = ioutil.WriteFile(fooPath, []byte{}, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(barPath, []byte{1, 2, 3}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Set up our repository - it's OK if we already registered it...
	f := NewFileRepository("file")
	repository.Register("file", f)

	// ...and some URIs - we know that they will not fail parsing
	foo := storage.NewFileURI(fooPath)
	bar := storage.NewFileURI(barPath)
	baz := storage.NewFileURI(bazPath)

	fooCanWrite, err := storage.CanWrite(foo)
	assert.True(t, fooCanWrite)
	assert.Nil(t, err)

	barCanWrite, err := storage.CanWrite(bar)
	assert.True(t, barCanWrite)
	assert.Nil(t, err)

	bazCanWrite, err := storage.CanWrite(baz)
	assert.True(t, bazCanWrite)
	assert.Nil(t, err)
}

func TestFileRepositoryParent(t *testing.T) {
	// Set up our repository - it's OK if we already registered it.
	f := NewFileRepository("file")
	repository.Register("file", f)

	// note the trailing slashes are significant, as they tend to belie a
	// directory

	parent, err := storage.Parent(storage.NewURI("file:///foo/bar/baz"))
	assert.Nil(t, err)
	assert.Equal(t, "file:///foo/bar/", parent.String())

	parent, err = storage.Parent(storage.NewFileURI("/foo/bar/baz/"))
	assert.Nil(t, err)
	assert.Equal(t, "file:///foo/bar/", parent.String())

	parent, err = storage.Parent(storage.NewURI("file://C:/foo/bar/baz/"))
	assert.Nil(t, err)
	assert.Equal(t, "file://C:/foo/bar/", parent.String())

	if runtime.GOOS == "windows" {
		// Only the Windows version of filepath will know how to handle
		// backslashes.
		uri := storage.NewURI("file://C:\\foo\\bar\\baz\\")
		assert.Equal(t, "file://C:/foo/bar/baz/", uri.String())
		uri = storage.NewFileURI("C:\\foo\\bar\\baz\\")
		assert.Equal(t, "file://C:/foo/bar/baz/", uri.String())

		parent, err = storage.Parent(uri)
		assert.Nil(t, err)
		assert.Equal(t, "file://C:/foo/bar/", parent.String())
	}

	_, err = storage.Parent(storage.NewURI("file:///"))
	assert.Equal(t, storage.URIRootError, err)

	if runtime.GOOS == "windows" {
		// This is only an error under Windows, on *NIX this is a
		// relative path to a directory named "C:", which is completely
		// valid.

		// This should cause an error, since this is a Windows-style
		// path and thus we can't get the parent of a drive letter.
		_, err = storage.Parent(storage.NewURI("file://C:/"))
		assert.Equal(t, storage.URIRootError, err)
	}

	// Windows supports UNIX-style paths. /C:/ is also a valid path.
	parent, err = storage.Parent(storage.NewURI("file:///C:/"))
	assert.Nil(t, err)
	assert.Equal(t, "file:///", parent.String())
}

func TestFileRepositoryChild(t *testing.T) {
	// Set up our repository - it's OK if we already registered it.
	f := NewFileRepository("file")
	repository.Register("file", f)

	p, _ := storage.Child(storage.NewURI("file:///foo/bar"), "baz")
	assert.Equal(t, "file:///foo/bar/baz", p.String())

	p, _ = storage.Child(storage.NewURI("file:///foo/bar/"), "baz")
	assert.Equal(t, "file:///foo/bar/baz", p.String())

	if runtime.GOOS == "windows" {
		// Only the Windows version of filepath will know how to handle
		// backslashes.
		uri := storage.NewURI("file://C:\\foo\\bar\\")
		assert.Equal(t, "file://C:/foo/bar/", uri.String())

		p, _ = storage.Child(uri, "baz")
		assert.Equal(t, "file://C:/foo/bar/baz", p.String())
	}
}

func TestFileRepositoryCopy(t *testing.T) {
	// Set up a temporary directory.
	dir, err := ioutil.TempDir("", "FyneInternalRepositoryFileTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create some files to test with.
	fooPath := path.Join(dir, "foo")
	barPath := path.Join(dir, "bar")
	err = ioutil.WriteFile(fooPath, []byte{1, 2, 3, 4, 5}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	foo := storage.NewFileURI(fooPath)
	bar := storage.NewFileURI(barPath)

	err = storage.Copy(foo, bar)
	assert.Nil(t, err)

	fooData, err := ioutil.ReadFile(fooPath)
	assert.Nil(t, err)

	barData, err := ioutil.ReadFile(barPath)
	assert.Nil(t, err)

	assert.Equal(t, fooData, barData)
}

func TestFileRepositoryMove(t *testing.T) {
	// Set up a temporary directory.
	dir, err := ioutil.TempDir("", "FyneInternalRepositoryFileTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create some files to test with.
	fooPath := path.Join(dir, "foo")
	barPath := path.Join(dir, "bar")
	err = ioutil.WriteFile(fooPath, []byte{1, 2, 3, 4, 5}, 0755)
	if err != nil {
		t.Fatal(err)
	}

	foo := storage.NewFileURI(fooPath)
	bar := storage.NewFileURI(barPath)

	err = storage.Move(foo, bar)
	assert.Nil(t, err)

	barData, err := ioutil.ReadFile(barPath)
	assert.Nil(t, err)

	assert.Equal(t, []byte{1, 2, 3, 4, 5}, barData)

	// Make sure that the source doesn't exist anymore.
	ex, err := storage.Exists(foo)
	assert.Nil(t, err)
	assert.False(t, ex)
}

func TestFileRepositoryListing(t *testing.T) {
	// Set up a temporary directory.
	dir, err := ioutil.TempDir("", "FyneInternalRepositoryFileTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create some files to tests with.
	fooPath := path.Join(dir, "foo")
	os.Mkdir(fooPath, 0755)
	os.Mkdir(path.Join(fooPath, "bar"), 0755)
	os.Mkdir(path.Join(fooPath, "baz"), 0755)
	os.Mkdir(path.Join(fooPath, "baz", "quux"), 0755)

	foo := storage.NewFileURI(fooPath)

	canList, err := storage.CanList(foo)
	assert.Nil(t, err)
	assert.True(t, canList)

	listing, err := storage.List(foo)
	assert.Nil(t, err)
	stringListing := []string{}
	for _, u := range listing {
		stringListing = append(stringListing, u.String())
	}
	assert.ElementsMatch(t, []string{"file://" + dir + "/foo/bar", "file://" + dir + "/foo/baz"}, stringListing)
}
