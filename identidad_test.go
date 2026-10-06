package main

import "testing"

// Identificadores INVENTADOS: nunca un NSA real en un test.
func TestIdentidadIdaYVuelta(t *testing.T) {
	const pid, nsa = uint64(1800000099), uint64(0xC0FFEE0012345678)
	lierIdentite(pid, nsa)
	if pidPublico(pid) != nsa || pidInterno(nsa) != pid {
		t.Fatal("la pareja PID<->NSA no se traduce en los dos sentidos")
	}
	if pidPublico(nsa) != nsa || pidInterno(pid) != pid {
		t.Fatal("las traducciones deben ser idempotentes")
	}
	for _, p := range []uint64{0, 2} {
		if pidPublico(p) != p {
			t.Fatalf("un PID de servicio (%d) no se toca", p)
		}
	}
}
