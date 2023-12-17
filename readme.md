# Documentation Technique : Infrastructure Whanos Jenkins

## Aperçu

Cette documentation technique donne un aperçu de l'infrastructure Whanos Jenkins, y compris la structure du projet, les Dockerfiles pour divers langages de programmation, la configuration Jenkins et le déploiement Kubernetes.

### Structure du Projet

Le projet se compose des composants suivants :

- **Scripts Shell :**
  - `build.sh` : Script Shell pour la configuration de l'environnement de développement, la construction et l'exécution de Jenkins dans un conteneur Docker, et la configuration d'un registre Docker local.
  - `stop.sh` : Script Shell pour arrêter et nettoyer l'environnement de développement, y compris les conteneurs et les images Docker.

- **Dockerfiles :**
  - `images/befunge/Dockerfile.base` : Dockerfile de base pour le langage Befunge.
  - `images/befunge/Dockerfile.standalone` : Dockerfile autonome pour le langage Befunge.
  - `images/c/Dockerfile.base` : Dockerfile de base pour le langage C.
  - `images/c/Dockerfile.standalone` : Dockerfile autonome pour le langage C.
  - `images/java/Dockerfile.base` : Dockerfile de base pour le langage Java.
  - `images/java/Dockerfile.standalone` : Dockerfile autonome pour le langage Java.
  - `images/javascript/Dockerfile.base` : Dockerfile de base pour le langage JavaScript.
  - `images/javascript/Dockerfile.standalone` : Dockerfile autonome pour le langage JavaScript.
  - `images/python/Dockerfile.base` : Dockerfile de base pour le langage Python.
  - `images/python/Dockerfile.standalone` : Dockerfile autonome pour le langage Python.
  - `jenkins/Dockerfile` : Dockerfile pour la configuration de Jenkins avec des configurations nécessaires.

- **Configuration Jenkins :**
  - `jenkins/jenkins.yml` : Fichier de configuration YAML pour Jenkins, comprenant le message système, les paramètres de sécurité et les configurations de plugins.
  - `jenkins/job_dsl.groovy` : Script Groovy pour créer dynamiquement des emplois Jenkins, comprenant des emplois pour la construction d'images de base et la liaison de projets.

- **Scripts Jenkins :**
  - `jenkins/findTech.sh` : Script Shell pour détecter le langage de programmation d'un projet en fonction de sa structure.
  - `jenkins/whanos.sh` : Script Shell pour containeriser des projets, construire des images Docker et les déployer dans Kubernetes.

- **Kubernetes :**
  - `kubernetes/go.mod` : Fichier de module Go pour le serveur API Kubernetes.
  - `kubernetes/main.go` : Code Go pour un serveur API simple qui déploie des images Docker en tant que déploiements Kubernetes.

## Utilisation

### Configuration de l'Environnement de Développement

1. Exécutez `build.sh` pour configurer l'environnement de développement. Ce script démarre Minikube, configure Docker, construit Jenkins et initialise un registre Docker local.

```bash
./build.sh
```

### Arrêt de l'Environnement de Développement

1. Exécutez `stop.sh` pour arrêter et nettoyer l'environnement de développement. Ce script supprime les conteneurs et les images Docker, arrête Minikube et tue les processus utilisant des ports spécifiques.

```bash
./stop.sh
```

### Configuration Jenkins

- Jenkins est configuré en utilisant le fichier `jenkins/jenkins.yml`. Il désactive l'assistant de configuration, définit des rôles et des permissions, et spécifie l'emplacement de la configuration Jenkins en tant que variable d'environnement.

- Les emplois Jenkins sont créés dynamiquement à l'aide du script `jenkins/job_dsl.groovy`. Il définit des emplois pour construire des images de base et lier des projets.

### Détection du Langage du Projet

- Le script `jenkins/findTech.sh` détecte le langage de programmation d'un projet en vérifiant des fichiers spécifiques (par exemple, `Makefile`, `package.json`, etc.).

### Containerisation et Déploiement des Projets

- Le script `jenkins/whanos.sh` containerise des projets en fonction de leur langage de programmation, construit des images Docker et les déploie dans Kubernetes.

### Déploiement Kubernetes et Forwarding de Port

Lorsque vous avez déployé un projet dans Kubernetes, n'oubliez pas de faire un *forward* de port pour accéder à l'application. Utilisez la commande suivante comme exemple, en remplaçant `pod/ramos-5b78c5cdd-8qbvb` par le nom réel du pod et `8082:3000` par le port local que vous souhaitez utiliser.

```bash
kubectl port-forward pod/ramos-5b78c5cdd-8qbvb 8082:3000
```

Cette commande permet de rediriger le trafic du port 8082 de votre machine locale vers le port 3000 du pod Kubernetes. Assurez-vous d'ajuster les valeurs en conséquence pour correspondre à votre configuration.

## Conclusion

Cette documentation technique donne un aperçu de l'infrastructure Whanos Jenkins, couvrant la structure du projet, les Dockerfiles, la configuration Jenkins et le déploiement Kubernetes. Elle sert de guide pour la configuration de l'environnement de développement, la gestion des emplois Jenkins et le déploiement de projets dans Kubernetes.