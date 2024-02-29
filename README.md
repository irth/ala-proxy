# ala-proxy

A proxy for when you don't want to upgrade your Arch installation as often as
you should.

When you `pacman -S` a package after a long time without updates, it might no
longer be available in the repositories. Using `-Sy` could result in an
inconsistent system state, and you might not want to do a full upgrade at the
time.

If that happens, this proxy will instead serve the old version of the package
from the [Arch Linux Archive](https://archive.archlinux.org/).

## Usage

```
~$ ala-proxy -h
Usage of ala-proxy:
  -archive string
        archive url (default "https://archive.archlinux.org")
  -listen string
        listen address (default ":8080")
  -meow
        meow
  -upstream string
        upstream repo url (default "https://arch.sakamoto.pl/$repo/os/$arch")
```

Run the proxy, set the first line of `/etc/pacman.d/mirrorlist` to:

```
Server = http://<address of this proxy>/$repo/os/$arch
```

Then, whenever a package is missing from the upstream repo, the proxy will try
to serve it from the Arch Linux Archive. If it exists on the upstream repo, the
proxy will just... proxy it.


## TODO

- [ ] GH actions
- [ ] systemd service
- [ ] PKGBUILD
- [ ] AUR package