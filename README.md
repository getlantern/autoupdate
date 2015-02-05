# Lantern Autoupdate

The `autoupdate` package provides lantern with the ability to request, download
and apply software patches over the network.

## General flow

![lanternautoupdates - general client](https://cloud.githubusercontent.com/assets/385670/6068159/a70eae50-ad3f-11e4-9cf6-42ed98e8f266.png)

At some point on the Lantern application's lifetime, an independent process
will be created, this process will periodically send local information (using a
proxy, if available) to an update server that will compare client's data
against a list of releases. When applicable, the server will generate a binary
patch and send a reply to the client containing the URL of the appropriate
patch. The client will download and apply the patch to itself so the new
version is ready at the next Lantern startup.

### Update server

![lanternautoupdates - server process](https://cloud.githubusercontent.com/assets/385670/6068177/cdaed85a-ad3f-11e4-802d-96c4cefe084e.png)

The update server holds a list of releases and waits for queries from clients.
Clients will send their own checksum and the server will compare that checksum
against the checksum of the latest release, if they don't match a binary diff
will be generated. This binary diff can be used by the client to patch itself.

### Download server

The update server may or may not be used as a download server. Clients will
pull binary diffs from this location, the actual patch's URL will be provided
by the update server. There is no need for using HTTPs on this download server,
as the binary diff's signatures will be checked in the client side.

### Client

![lanternautoupdates - auto update process](https://cloud.githubusercontent.com/assets/385670/6068187/e34459ec-ad3f-11e4-9b9e-1874f05fa735.png)

A client is aware of its own checksum and will send it to an update server
periodically. When the update server replies with a special message meaning
that a new version is available, the client will download the binary patch,
check its signature and apply it to the program's path, replacing the old
binary.
