# Lantern Autoupdate

The `autoupdate` package provides [Lantern][1] with the ability to request,
download and apply software updates over the network with minimal interaction.
At this time, `autoupdate` relies on the [go-update][2] package and the
[equinox update server][3].

## General flow

![lanternautoupdates - general client](https://cloud.githubusercontent.com/assets/385670/6096628/a4b53b78-af61-11e4-829a-f3be4a011846.png)

At some point on the Lantern application's lifetime, an independent process
will be created, this process will periodically send local information (using a
proxy, if available) to an update server that will compare client's data
against a list of releases. When applicable, the server will generate a binary
patch and send a reply to the client containing the URL of the appropriate
patch. The client will download and apply the patch to its executable file so
the new version is ready the next time Lantern starts.

### Update server

![lanternautoupdates - server process](https://cloud.githubusercontent.com/assets/385670/6096630/a85e20e6-af61-11e4-85e0-03e2f0740057.png)

The update server holds a list of releases and waits for queries from clients.
Clients will send their own checksum and the server will compare that checksum
against the checksum of the latest release, if they don't match a binary diff
will be generated. This binary diff can be used by the client to patch itself.

### Download server

The update server may or may not be used as a download server. Clients will
pull binary diffs from this location, the actual patch's URL will be provided
by the update server.

### Client

![lanternautoupdates - auto update process](https://cloud.githubusercontent.com/assets/385670/6096629/a6ede304-af61-11e4-940c-c56a4c28b3d5.png)

A client will compute the checksum of its executable file and will send it to
an update server periodically. When the update server replies with a special
message meaning that a new version is available, the client will download the
binary patch, apply it to a temporary file and check the signature, if the
signature is what the client expects, the original executable will be replaced
with the patched one.

[1]: https://getlantern.org/
[2]: https://github.com/inconshreveable/go-update
[3]: https://equinox.io/
