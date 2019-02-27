# Secret seed: SANFNPZPA4LWBD3RPDSCJU63KCBU3OBFOM5FFBJCGIOCVIABMRTKBAU2
# Address: GCYLTPOU7IVYHHA3XKQF4YB4W4ZWHFERMOQ7K47IWANKNBFBNJJNEOG5
# Asset Code Issued: USD
account :source, Blocksafe::KeyPair.from_seed("SANFNPZPA4LWBD3RPDSCJU63KCBU3OBFOM5FFBJCGIOCVIABMRTKBAU2")

# Secret seed: SABP5P625YBETJV4BCEWQD674ED4FF4QVNBRL6TQCRODJUWBNJBMND5O
# Address: GCSX4PDUZP3BL522ZVMFXCEJ55NKEOHEMII7PSMJZNAAESJ444GSSJMO
# Asset Code Issued: EUR
account :dest, Blocksafe::KeyPair.from_seed("SABP5P625YBETJV4BCEWQD674ED4FF4QVNBRL6TQCRODJUWBNJBMND5O")

# Secret seed: SCCUFFUANIXJPAWBHDXZXY5D4GB32QPM6MOUWDD6PTYBLPE6JVYZFE76
# Address: GCFZWN3AOVFQM2BZTZX7P47WSI4QMGJC62LILPKODTNDLVKZZNA5BQJ3
# Asset Code Issued: USD
account :issuer, Blocksafe::KeyPair.from_seed("SCCUFFUANIXJPAWBHDXZXY5D4GB32QPM6MOUWDD6PTYBLPE6JVYZFE76")

# Secret seed: SAISD7SISIIW5YNQ7GY5727L6MOFS667K3LVIPYPPUBIPCRQUORFLQMN
# Address: GAB7GMQPJ5YY2E4UJMLNAZPDEUKPK4AAIPRXIZHKZGUIRC6FP2LAQSDN
# Asset Code Issued: USD
account :another, Blocksafe::KeyPair.from_seed("SAISD7SISIIW5YNQ7GY5727L6MOFS667K3LVIPYPPUBIPCRQUORFLQMN")

use_manual_close

# create accounts
create_account :source
create_account :dest
create_account :issuer
create_account :another
close_ledger

# add a trust line so the :source account lists USD:issuer as a currency it holds
trust :source, :issuer, "USD"
# add a trust line so the :dest account lists USD:another as a currency it holds
trust :dest, :another, "USD"
close_ledger