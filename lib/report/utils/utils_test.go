package utils

import (
	"testing"
)

func TestCharactersTheSame(t *testing.T) {

    if CharactersTheSame("http://www.google.com/jeffrey", "jeffrey") != 7 {
        t.Errorf("Doesn't work right")
    }

    if CharactersTheSame("http://www.google.com/jeffrey", "jeffey") != 4 {
        t.Errorf("Doesn't work right")
    }

    if CharactersTheSame("http://www.google.com/jeffcrapjeffrey", "jeffrey") != 7 {
        t.Errorf("Doesn't work right")
    }

    if CharactersTheSame("http://www.google.com/jesus", "jeffrey") != 2 {
        t.Errorf("Doesn't work right")
    }

}
