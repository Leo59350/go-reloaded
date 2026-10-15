# Go-reloaded

Programme Go utilisant uniquement la bibliothèque standard. Il lit un fichier,
transforme son texte et écrit le résultat dans un autre fichier.

## Lancer le programme

Ouvre un terminal dans ce dossier, puis lance :

```sh
go run . sample.txt result.txt
```

Le fichier `result.txt` contient le résultat. Attention : un fichier de sortie
existant est remplacé après une transformation réussie.

Pour lancer les tests :

```sh
go test ./...
```

## Les règles

| Balise | Effet | Exemple |
| --- | --- | --- |
| `(hex)` | Base 16 vers base 10 | `1E (hex)` vers `30` |
| `(bin)` | Base 2 vers base 10 | `10 (bin)` vers `2` |
| `(up)` | Majuscules | `hello (up)` vers `HELLO` |
| `(low)` | Minuscules | `HELLO (low)` vers `hello` |
| `(cap)` | Première lettre majuscule, reste minuscule | `hELLO (cap)` vers `Hello` |
| `(up, 2)` | Deux mots précédents | `hello world (up, 2)` vers `HELLO WORLD` |

Le nombre de mots fonctionne aussi avec `(low, N)` et `(cap, N)`.
Les balises appliquées disparaissent. La ponctuation ne compte pas comme un mot.
Les signes `. , ! ? : ;` se collent au mot précédent. Les groupes restent ensemble.
Les quotes simples se collent au contenu : `' hello world '` devient `'hello world'`.
Les apostrophes internes, comme dans `don't`, restent dans leur mot.
`a` et `A` deviennent `an` et `An` devant a, e, i, o, u ou h,
sans distinction de casse, conformément au sujet.
