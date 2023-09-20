<h1 align="center">Discord Spotify Status</h1>

<div align="center">
 
[![Stars](https://img.shields.io/github/stars/bylkamar/telegram-rss?style=social)](https://github.com/bylkamar/telegram-rss)
[![Forks](https://img.shields.io/github/forks/bylkamar/telegram-rss?style=social
)](https://github.com/bylkamar/telegram-rss)
[![Issues](https://img.shields.io/github/issues/bylkamar/telegram-rss
)](https://github.com/bylkamar/telegram-rss)

**Synchronisé votre écoute Spotify à votre profil Discord**


</div>

> **⚠️ L'utilisation de ce logiciel est uniquement à but éducatif, toutes actions menées sur votre compte sont entièrement sous votre responsabilité.** 

## 📦 Installation 
* Lancer directement la version compiler ou le build directement depuis le code source
```go
go build
```

* Configurer le fichier JSON
```json
{
    "session": "Votre Cookie de session spotify (sp_dc=XXXXXXX..)",
    "discord_token": "Token"
}
```
### Comment obtenir le cookie de session?

- Aller sur web.spotify.com
- Ouvrir l'outil de développement web
- Copier et coller la valeur du cookie "sp_dc"dans le fichier `session.json`


<img src="https://imgur.com/G1Vtkhd.png">

<br/>

## 🚀 Utilisation

> ➕ Executer le logiciel

> 🎉 Vous n'avez plus qu'a lancer une musique sur Spotify et il se synchronisera automatiquement sur votre status discord
## 📚 To Do List

* Version Graphique
* Récupération automatique du cookie de session
* Version RichPresence (RPC)
* Possibilité de modifier le format et le nombre de lignes à affiché

> Je suis ouvert <a href="https://github.com/bylkamar/telegram-rss/pulls">aux idées</a> ou bien si vous avez reperéz <a href="https://github.com/bylkamar/telegram-rss/issues">des bugs</a>, soumettez-les moi.



## 👥 Contributeurs

<p align="center">
  <a href="https://github.com/bylkamar/telegram-rss/graphs/contributors">
    <img src="https://contrib.rocks/image?repo=bylkamar/telegram-rss" />
  </a>
</p>

## License

**Apache-2.0 license**
