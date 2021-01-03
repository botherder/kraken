rule eicar_test {
    strings:
        $eicar = "$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!"

    condition:
        all of them
}
