{ pkgs
, dbname ? "goshrt"
, dbuser ? "goshrt"
, dbpass ? "trhsog"
, dbport ? "6000"
, pgdata ? ".devshell/db"
}:

with pkgs; [
  glibcLocales
  postgresql
  pgcli
  (writeScriptBin "pgnix-init" ''
    initdb -D ${pgdata} -U postgres
    pg_ctl -D ${pgdata} -l ${pgdata}/postgres.log  -o "-p ${dbport} -k /tmp -i" start
    createdb --port=${dbport} --host=localhost --username=postgres -O postgres ${dbname}
    psql -d postgres -U postgres -h localhost -p ${dbport} -c "create user ${dbuser} with encrypted password '${dbpass}';"
    psql -d postgres -U postgres -h localhost -p ${dbport} -c "grant all privileges on database ${dbname} to ${dbuser};"
  '')
  (writeScriptBin "pgnix-start" ''
    pg_ctl -D ${pgdata} -l ${pgdata}/postgres.log  -o "-p ${dbport} -k /tmp -i" start
  '')
  (writeScriptBin "pgnix-pgcli" ''
    PGPASSWORD=${dbpass} pgcli -h localhost -p 6000 -U goshrt
  '')
  (writeScriptBin "pgnix-psql" ''
    PGPASSWORD=${dbpass} psql -d ${dbname} -U ${dbuser} -h localhost -p ${dbport}
  '')
  (writeScriptBin "pgnix-status" ''
    pg_ctl -D ${pgdata} status
  '')
  (writeScriptBin "pgnix-restart" ''
    pg_ctl -D ${pgdata} restart
  '')
  (writeScriptBin "pgnix-stop" ''
    pg_ctl -D ${pgdata} stop
  '')
  (writeScriptBin "pgnix-purge" ''
    pg_ctl -D ${pgdata} stop
    rm -rf .devshell/db
  '')
]
