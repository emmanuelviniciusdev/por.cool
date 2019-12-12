<template>
  <div>
    <div class="date-filter">
      <b-field label="ano" class="input-select">
        <b-select placeholder="selecione um ano">
          <option value="2018">2018</option>
          <option value="2019">2019</option>
        </b-select>
      </b-field>
      <b-field label="mês" class="input-select">
        <b-select placeholder="selecione um mês">
          <option value="janeiro">janeiro</option>
          <option value="fevereiro">fevereiro</option>
        </b-select>
      </b-field>
    </div>
  </div>
</template>

<script>
export default {
    name: "FilterByDate"
};
</script>

<style lang="scss" scoped>
.input-select {
  display: inline-block;
  margin-right: 15px;
}
</style>