<template>
  <div>
    <div class="date-filter">
      <b-field label="ano" class="input-select">
        <div
          class="select"
          :class="{ 'is-loading': spendingDatesList.length === 0 }"
        >
          <select @change="onDateChange('yearInput', $event.target.value)">
            <option
              v-for="spendingDate in spendingDatesList"
              :key="spendingDate.year"
              :selected="spendingDate.year == currentYear"
              >{{ spendingDate.year }}</option
            >
          </select>
        </div>
      </b-field>
      <b-field label="mês" class="input-select">
        <div v-if="spendingDatesListOnlyMonths">
          <div
            class="select"
            :class="{ 'is-loading': spendingDatesList.length === 0 }"
          >
            <select @change="onDateChange('monthInput', $event.target.value)">
              <option
                v-for="month in spendingDatesListOnlyMonths"
                :key="month"
                :selected="month === currentMonth"
                >{{ month }}</option
              >
            </select>
          </div>
        </div>
      </b-field>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import moment from "moment";

// Helpers
import dateAndTime from "../helpers/dateAndTime";

export default {
  name: "FilterByDate",
  data() {
    return {
      currentYear: null,
      currentMonth: null
    };
  },
  computed: {
    ...mapState({
      userData: state => state.user.user,
      lookingAtSpendingDate: state => state.user.user.lookingAtSpendingDate,
      spendingDatesList: state => state.expenses.spendingDatesList
    }),
    spendingDatesListOnlyMonths() {
      if (this.spendingDatesList.length > 0) {
        const filteredSpendingDatesList = this.spendingDatesList.filter(
          date => date.year === parseInt(this.currentYear)
        );

        if (filteredSpendingDatesList[0] !== undefined)
          return filteredSpendingDatesList[0].months;
      }
    }
  },
  watch: {
    lookingAtSpendingDate() {
      if (this.lookingAtSpendingDate !== undefined) this.setCurrentDate();
    }
  },
  methods: {
    setCurrentDate() {
      this.$store.dispatch("expenses/setSpendingDatesList", {
        userUid: this.userData.uid,
        lookingAtSpendingDate: this.lookingAtSpendingDate
      });

      const formatedLookingAtSpendingDate = dateAndTime.extractOnly(
        this.lookingAtSpendingDate,
        ["year", "month"]
      );

      this.currentYear = formatedLookingAtSpendingDate.year;
      this.currentMonth = formatedLookingAtSpendingDate.month;
    },
    onDateChange(from, value) {
      this.currentYear = from === "yearInput" ? value : this.currentYear;
      this.currentMonth = from === "monthInput" ? value : "janeiro";

      const monthInNumber = dateAndTime.months.indexOf(this.currentMonth) + 1;

      const dateFormat = `${this.currentYear}-${monthInNumber}-01`;
      const spendingDate = moment(dateFormat, "YYYY-MM-DD")
        .startOf("days")
        .toDate();

      this.$emit("onDateChange", spendingDate);
    }
  },
  created() {
    this.setCurrentDate();
  }
};
</script>

<style lang="scss" scoped>
.input-select {
  display: inline-block;
  margin-right: 15px;
}
</style>
