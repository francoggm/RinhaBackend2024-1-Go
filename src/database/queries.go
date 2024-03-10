package database

const CreateUsersQuery = `
	MERGE (:Usuario {id: 1, limite: 100000, saldo: 0})
	MERGE (:Usuario {id: 2, limite: 80000, saldo: 0})
	MERGE (:Usuario {id: 3, limite: 1000000, saldo: 0})
	MERGE (:Usuario {id: 4, limite: 10000000, saldo: 0})
	MERGE (:Usuario {id: 5, limite: 500000, saldo: 0})
`

const ExtractQuery = `
	MATCH (u:Usuario {id: $id})
	OPTIONAL MATCH (u)-[:REALIZOU]->(t)

	WITH u, t
	ORDER BY t.data DESC LIMIT 10

	WITH u,
			CASE 
					WHEN t IS NOT NULL THEN
							{tipo: t.tipo, valor: abs(t.valor), descricao: t.descricao, data: t.data}
					ELSE
							NULL
			END as ts

	RETURN u.saldo AS saldo, u.limite AS limite, collect(ts) AS transacoes`

const TransactionQuery = `
	MATCH (u:Usuario {id: $id})

	SET u._LOCK_ = true

	WITH u, 
		CASE 
			WHEN u IS NOT NULL AND ($tipo = 'd' AND u.saldo + $valor < -1 * u.limite) THEN 
				NULL
			ELSE 
				{valor: $valor, tipo: $tipo, descricao: $descricao, data: timestamp()} 
		END as t

	FOREACH(_ IN CASE WHEN t IS NULL THEN [] ELSE [1] END |
		CREATE (u)-[:REALIZOU]->(nt:Transacao)
		SET nt = t,
				u.saldo = u.saldo + t.valor
	)

	REMOVE u._LOCK_

	RETURN u.saldo AS saldo, u.limite AS limite, t as transacao`
