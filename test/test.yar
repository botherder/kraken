rule TEST_RULE {
    strings:
        $a = "This is a test Yara rule"

    condition:
        $a
}
