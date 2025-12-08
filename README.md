# por.cool

## Project setup

### 1. Install dependencies
```
yarn install
```

### 2. Configure Firebase
Copy the `environment.example.js` file and create your environment configuration files:

```bash
cp environment.example.js environment.development.js
cp environment.example.js environment.production.js
```

Then edit both files with your Firebase project credentials:
- `environment.development.js` - for development/testing
- `environment.production.js` - for production

### Compiles and hot-reloads for development
```
yarn serve
```

### Compiles and minifies for production
```
yarn build
```

### Run your unit tests
```
yarn test:unit
```

### Lints and fixes files
```
yarn lint
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).
