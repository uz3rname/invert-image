package db

type Store interface {
  GetLastImages(n int) []*ImagePair
  AddImage(orig, neg, origMime, negMime, hash string) *ImagePair
  FindImageByHash(hash string) (*ImagePair, bool)
}
