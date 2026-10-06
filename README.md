# Just Shapes & Beats on Nextendo

Original NEX matchmaking service for Just Shapes & Beats. The game uses Pia for
peer-to-peer play; this server provides authentication, session discovery and NAT
traversal. It does not relay traffic to another game's live service.

The game's measured server ID is `0x2a699600`, and the access key
`1b2ee2f7` was verified against a captured PRUDP CONNECT signature. The code
targets NEX 4.6.5. Juan confirmed a complete multiplayer game between two
consoles on 2026-10-06. A player on another network then joined a session with
successful direct Pia traversal. Public matchmaking and online multiplayer are
deployed; unknown RMC methods remain logged by protocol and method number.

The decisive fix was matching the station identity and join response settings
used by Mario Strikers' newer Pia: `PreservePiaStationIdentity`,
`JoinRespExistingCount`, `SessionPartPersists`, and the type-8 keepalive. Before
that, consoles found the same lobby and could complete NAT traversal but left
the Pia session. The secure station also uses `prudps`, and the deployment
mounts the same NNCS observation file as SMM2 and SSBU.

The JSAB binary contains Pia jobs for relay negotiation
(`PrepareNatTraversalByRelay`, `SendRelayConnectionRequest`, and
`RelayRouteManageJob`). The current pair relay only forwards UDP; the NEX
`GetRelaySignatureKey` handler advertises no relay. Its packet forwarding
therefore does not prove that JSAB's Pia has accepted a relay connection. The
successful console tests used direct Pia traversal with this relay off.

Run on the VPS with `deploy-vps.sh` after placing the compiled `jsab-server` in
`/home/juan/jsab/`. The script checks that its ports and container name are free.
The auth host is routed through Traefik; secure traffic uses TCP port 60024.
The deployment mounts `/opt/mk8nex-nat` read-only for NNCS UDP observations.
The dashboard is available only on the Docker `coolify` network at
`http://jsab:8113/api/stats`.

Direct Pia connections can fail for some NAT combinations. The UDP relay is
off by default. To test one measured pair, write the VPS public IP to
`/home/juan/jsab/relay_on` and both internal Nextendo PIDs, one per line, to
`/home/juan/jsab/relay_allow`. Prefix a PID with `!` to force relay for a
test pair without waiting for two direct failures. Both files are reloaded
every two seconds. Remove `relay_on` to turn the relay off. Only ports
39000–39015/UDP are exposed for this service.

Build on Juan's computer with limited resources:

```sh
GOFLAGS=-p=4 GOMAXPROCS=4 GOMEMLIMIT=2GiB go build -o ./jsab-server ./...
```

The service uses Nextendo's own NEX code and is licensed under PolyForm Shield
1.0.0; see `LICENSE.md`.
