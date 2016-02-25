package svfs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/ncw/swift"
	"golang.org/x/net/context"
)

type Object struct {
	name      string
	path      string
	so        *swift.Object
	c         *swift.Container
	cs        *swift.Container
	p         *Directory
	segmented bool
}

func (o *Object) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Size = o.size()
	a.Mode = 0600
	a.Mtime = o.so.LastModified
	a.Ctime = o.so.LastModified
	a.Crtime = o.so.LastModified
	return nil
}

func (o *Object) Export() fuse.Dirent {
	return fuse.Dirent{
		Name: o.Name(),
		Type: fuse.DT_File,
	}
}

func (o *Object) open(mode fuse.OpenFlags) (oh *ObjectHandle, err error) {
	oh = &ObjectHandle{t: o}

	// Append mode is not supported
	if mode&fuse.OpenAppend == fuse.OpenAppend {
		return nil, fuse.ENOTSUP
	}
	if mode.IsReadOnly() {
		oh.rd, _, err = SwiftConnection.ObjectOpen(o.c.Name, o.so.Name, false, nil)
		return oh, err
	}
	if mode.IsWriteOnly() {
		oh.w, err = SwiftConnection.ObjectCreate(o.c.Name, o.so.Name, false, "", "application/octet-sream", nil)
		return oh, err
	}

	return nil, fuse.ENOTSUP
}

func (o *Object) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	return o.open(req.Flags)
}

func (o *Object) Name() string {
	return o.name
}

func (o *Object) size() uint64 {
	return uint64(o.so.Bytes)
}

var (
	_ Node          = (*Object)(nil)
	_ fs.Node       = (*Object)(nil)
	_ fs.NodeOpener = (*Object)(nil)
)
