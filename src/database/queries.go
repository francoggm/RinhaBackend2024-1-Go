package database

const ExtractQuery = `
		MATCH (u:Usuario {id: $id})-[:REALIZOU]->(t)
		WITH u, t
		ORDER BY t.date DESC LIMIT 10

		WITH u, collect({tipo: t.tipo, valor: abs(t.valor), descricao: t.descricao, data: t.data}) AS transacoes
		RETURN u.saldo AS saldo, u.limite AS limite, transacoes
	`

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

	SET u._LOCK_ = false

	RETURN u.saldo AS saldo, u.limite AS limite, t as transacao`
