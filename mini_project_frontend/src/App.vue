<script setup>
import { ref } from 'vue'
import { loadScript } from "vue-plugin-load-script";



const options = {
  region: 'US',
  resolution: 'provinces'
};
const queried_states_with_prefix = ref([]);
var drawRegionsMap;
loadScript('https://www.gstatic.com/charts/loader.js')
  .then(() => {
    google.charts.load('current', {
        'packages':['geochart'],
      });
      drawRegionsMap = () => {
        var data_table = []
        if (queried_states_with_prefix.value.length == 1) {
          data_table = [[queried_states_with_prefix.value[0]]];
        }
        var data = google.visualization.arrayToDataTable([
          ['State'],
        ].concat(data_table));

        var options = {
          region: 'US',
          resolution: 'provinces'
        };

        var chart = new google.visualization.GeoChart(document.getElementById('regions_div'));

        chart.draw(data, options);
      }
      google.charts.setOnLoadCallback(drawRegionsMap);
  })
  .catch((e) => {
    console.log(e)
  })


async function onInput(e) {
  queried_states_with_prefix.value = []
  fetch(`http://localhost:3000/query?state_prefix=${e.target.value}`)
    .then((res) => {
      return res.json();
    })
    .then((states_with_prefix) => {
      queried_states_with_prefix.value = states_with_prefix;
      drawRegionsMap();
    })
    .catch((err) => {
      console.log('Error: ', err)
    })
}

</script>

<template>
  <component src="https://www.gstatic.com/charts/loader.js" :is="'script'"></component>
  <div id="regions_div" style="width: 900px; height: 500px;"></div>  
  <!-- <input :value="state_prefix" @input="onInput" placeholder="Type here">
  <p v-if="queried_states_with_prefix">{{ queried_states_with_prefix }}</p> -->
  <input list="States" @input="onInput"/>
  
  <datalist id="States">
      <option v-for="state in queried_states_with_prefix" :value="state">{{state}}</option> -->
  </datalist>
  
  <div id="regions_div" style="width: 900px; height: 500px;"></div>
</template>