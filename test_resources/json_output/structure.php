<?php

$json = stream_get_contents(STDIN);
$json = preg_replace("/: -?(?:1.#IN[DF]|nan|inf)/", ": null", $json);
$data = json_decode($json, true);
if ($data === null) {
    die("JSON Decode Error!\n");
}

$lines = [];

foreach ($data as $image) {
    foreach ($image['image'] as $key => $val) {

        $varname = ucfirst($key);

        $type = "string";
        if (is_bool($val)) {
            $type = "bool";
        } else if (is_int($val)) {
            $type = "int64";
        } else if (is_float($val) || is_double($val)) {
            $type = "float64";
        } else if (is_array($val)) {
            $type = "*ImageMagick$varname";
        }

        $tags = sprintf('`json:"%s"`', $key);

        $lines[] = sprintf("%-17s %-60s %s", $varname, $type, $tags);

    }
}

sort($lines);
$lines = array_unique($lines);
echo implode("\n", $lines);
