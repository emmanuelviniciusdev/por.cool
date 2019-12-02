<template>
  <div id="app">
    <h1>por.<span>cool</span></h1>
    <img src="assets/images/pig.png" alt="porcool">
    <p>você <b>não perde</b> por esperar</p>
  </div>
</template>

<script>
export default {
  name: 'app'
};
</script>

<style>
* {
  margin: 0;
  padding: 0;
  font-family: 'Lora', serif;
}

#app {
  width: 100%;
  /* background: #ccc; */
  overflow: auto;
  margin: 0 auto;
}

h1 {
  font-size: 40px;
  margin: 20px;
  display: block;
}
h1 span {
  color: #1BB450;
}

img {
  width: 270px;
  margin: 0 auto;
  margin-top: 30px;
  display: block;
}

p {
  text-align: center;
  margin-top: 40px;
  font-size: 30px;
}

@media only screen and (min-width: 1000px) {
  #app {width: 80%;}
  p {font-size: 35px;}
}
@media only screen and (min-width: 1200px) {
  #app {width: 1000px;}
}
</style>
