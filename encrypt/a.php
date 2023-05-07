<?php

$str = 'sky147258369';
//echo password_hash($str, PASSWORD_DEFAULT);
$hash = '$2y$10$KSvTmoYhv29RVmur2I83E.Wd6QQ5nS/TcwWZyS/W.Sw.NygCozFDuroot@bd92912795f4';
if (password_verify($str, $hash)) {
    echo 'Password is valid!';
} else {
    echo 'Invalid password.';
}
