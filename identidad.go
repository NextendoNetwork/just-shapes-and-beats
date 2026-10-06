package main

// Identidad pública: la consola debe ver su NSA como PID, no nuestro PID interno.
//
// Misma corrección que smm2-kazu/smm2_identite_publique.go (salas de amigos de SMM2, donde el
// anfitrión salía invisible). Aquí es una HIPÓTESIS (2026-10-06): JSAB usa Pia 5.37, que desde la
// 5.27 envía un "identification token" en la conexión entre consolas. Si ese token es la identidad
// de la consola y el anfitrión lo compara con los participantes que le dio NEX (PID internos), lo
// rechaza: hole-punch correcto y luego el visitante abandona (0x32.1, error 2618-0513/0502).
//
// Dentro del servidor nada cambia: conexiones y sesiones siguen con el PID interno. La traducción
// va en la frontera, en nextendo-nex (StreamOut.PID / StreamIn.PID, StationURL, Param2).
//
// Solo en memoria: /app es de solo lectura en el contenedor. Basta para JSAB, porque los jugadores
// que importan están conectados y su pareja NSA<->PID se aprende en su propio login.
//
// Interruptor: JSAB_IDENTIDAD_OFF=1 en el entorno.

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"

	nex "github.com/NextendoNetwork/nextendo-nex"
)

// premierPIDCompte: nuestros PID de cuenta empiezan aquí y caben en 32 bits. Por debajo son PID de
// servicio (el servidor seguro, 2): no se tocan.
const premierPIDCompte = 1_800_000_000

var (
	identMu  sync.Mutex
	nsaDePID = map[uint64]uint64{}
	pidDeNSA = map[uint64]uint64{}
)

func lierIdentite(pid, nsa uint64) {
	if pid < premierPIDCompte || pid > math.MaxUint32 || nsa <= math.MaxUint32 {
		return
	}
	identMu.Lock()
	defer identMu.Unlock()
	if viejo, ok := pidDeNSA[nsa]; ok && viejo != pid {
		delete(nsaDePID, viejo)
	}
	if viejo, ok := nsaDePID[pid]; ok && viejo != nsa {
		delete(pidDeNSA, viejo)
	}
	nsaDePID[pid] = nsa
	pidDeNSA[nsa] = pid
}

// pidPublico: lo que la consola debe ver para este PID. Nunca bloquea.
func pidPublico(pid uint64) uint64 {
	if pid < premierPIDCompte || pid > math.MaxUint32 {
		return pid
	}
	identMu.Lock()
	nsa, ok := nsaDePID[pid]
	identMu.Unlock()
	if ok {
		return nsa
	}
	// El jugador autenticado: su NSA es el nombre que dio al Auth.
	if nom, ok := nex.LoginNameForPID(pid); ok {
		if n, err := strconv.ParseUint(nom, 10, 64); err == nil && n > math.MaxUint32 {
			lierIdentite(pid, n)
			return n
		}
	}
	return pid
}

// pidInterno: nuestro PID para un identificador que viene de la consola.
func pidInterno(id uint64) uint64 {
	if id <= math.MaxUint32 {
		return id
	}
	identMu.Lock()
	p, ok := pidDeNSA[id]
	identMu.Unlock()
	if ok {
		return p
	}
	if p, st := resolveNSAtoPID(id); st == nsaOK {
		lierIdentite(p, id)
		return p
	}
	return id
}

// instalarIdentidad engancha la traducción en los ajustes de Auth y del servidor seguro.
func instalarIdentidad(reglajes ...*nex.Settings) {
	if os.Getenv("JSAB_IDENTIDAD_OFF") == "1" {
		fmt.Println("[JSAB identidad] DESACTIVADA: la consola ve el PID interno")
		return
	}
	for _, s := range reglajes {
		s.PIDPublic = pidPublico
		s.PIDInterne = pidInterno
	}
	fmt.Println("[JSAB identidad] activa: la consola ve su NSA como PID")
}
