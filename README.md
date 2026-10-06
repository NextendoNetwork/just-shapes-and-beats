# Just Shapes & Beats on Nextendo

Original NEX matchmaking service for Just Shapes & Beats. The game uses Pia for
peer-to-peer play; this server provides authentication, session discovery and NAT
traversal. It does not relay traffic to another game's live service.

The game's measured server ID is `0x2a699600`, and the access key
`1b2ee2f7` was verified against a captured PRUDP CONNECT signature. The code
targets NEX 4.6.5. Console tests confirmed public matchmaking, but joining a
second player returned 2618-0513. An opt-in relay received traffic from both
players and forwarded packets in both directions without resolving that error.
The secure station was corrected from `prudp` to `prudps`, and ticketed secure
CONNECTs are confirmed. A direct-connection test still failed because the host's
public station carried its TCP port instead of the UDP port in `ReplaceURL`.
The per-title UDP-port fallback is now deployed for a new console test. Do not
list JSAB as playable until a multiplayer match starts. Unknown RMC methods
are logged by protocol and method number.

Run on the VPS with `deploy-vps.sh` after placing the compiled `jsab-server` in
`/home/juan/jsab/`. The script checks that its ports and container name are free.
The auth host is routed through Traefik; secure traffic uses TCP port 60024.
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
GOFLAGS=-p=4 GOMAXPROCS=4 GOMEMLIMIT=2GiB go build -o /home/juanjo/Nextendo/just-shapes-and-beats/jsab-server ./...
```

The service uses Nextendo's own NEX code and is licensed under PolyForm Shield
1.0.0; see `LICENSE.md`.
