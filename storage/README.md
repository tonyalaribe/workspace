# storage
--
    import "gitlab.com/middlefront/workspace/stoRAGE"

Package storage makes it possible to perform file related operation from the
rest of workspace without worrying about the underlying storage implementation.

The storage mechanism to use, and the relevant connection details should be
specified in the workspace.yaml file.

.workspace.yaml

The .workspcae.yaml file coul have the following storage related files:

- persistence-type: This represents the storage backend to be used. eg local,
s3, openstack,

## Usage

#### func  GetByPath

```go
func GetByPath(path string) (stow.Item, error)
```
GetByPath is able to load a flle from stowBucket given the file ID(which is
usually a path to the file on the storage bucket)

#### func  GetByURL

```go
func GetByURL(urlstr string) (stow.Item, error)
```
GetByURL is able to load a file from stowLoc (stow location) which is bucket
independent, given a url which usually encodes information like the storage
backend kind, and the bucket the files exist in.

#### func  StorageInit

```go
func StorageInit() error
```
StorageInit initilizes stow with the storage details from the .workspace.yaml
configuration. StorageInit creates a stowBucket and stowLoc instance, which can
be shared during the lifetime of the application

#### type StowFile

```go
type StowFile struct {
	ID   string
	Name string
	URL  string
	Size int64
}
```


#### func  UploadBase64

```go
func UploadBase64(path string, b64 string) (StowFile, error)
```
UploadBase64 recieves a base64 string, which it converts to the underlying file
and uploads via UploadStream

#### func  UploadStream

```go
func UploadStream(path string, r io.Reader, size int64) (StowFile, error)
```
UploadStream recieves an io.reader which is used stream the file to the
repective storage backend. This is ideal, as it prevents sitting the file(bytes)
in memory during uploads. Path represents they path to the file, which might be
a single filename, or follow a backslack delimited structure.
