# Maintainer: Atila de Freitas <atiladefreitas@users.noreply.github.com>
pkgname=tical
pkgver=0.1.0
pkgrel=1
pkgdesc="A minimalist Tokyo Night themed terminal calculator (TUI)"
arch=('x86_64')
url="https://github.com/atiladefreitas/tical"
license=('MIT')
makedepends=('go')
source=("$pkgname-$pkgver.tar.gz::$url/archive/v$pkgver.tar.gz")
sha256sums=('SKIP') # run `updpkgsums` once the v0.1.0 release tarball exists

prepare() {
	cd "$pkgname-$pkgver"
	export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
	go mod download
}

build() {
	cd "$pkgname-$pkgver"
	export CGO_CPPFLAGS="${CPPFLAGS}"
	export CGO_CFLAGS="${CFLAGS}"
	export CGO_CXXFLAGS="${CXXFLAGS}"
	export CGO_LDFLAGS="${LDFLAGS}"
	export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"
	go build -o "$pkgname" .
}

check() {
	cd "$pkgname-$pkgver"
	go test ./...
}

package() {
	cd "$pkgname-$pkgver"
	install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
	install -Dm644 LICENSE "$pkgdir/usr/share/licenses/$pkgname/LICENSE"
	install -Dm644 README.md "$pkgdir/usr/share/doc/$pkgname/README.md"
}
