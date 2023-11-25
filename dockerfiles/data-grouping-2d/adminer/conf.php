<?php

define('ADMINER_DIR', '/usr/share/adminer');

function adminer_object() {
    include_once "/usr/share/adminer/plugins/plugin.php";
    include_once "/usr/share/adminer/plugins/login-password-less.php";
    return new AdminerPlugin(array(
        // TODO: inline the result of password_hash() so that the password is not visible in source codes
        new AdminerLoginPasswordLess(password_hash("changeme", PASSWORD_DEFAULT)),
    ));
}

include ADMINER_DIR . "/adminer.php";
?>
