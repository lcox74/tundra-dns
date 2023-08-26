<script lang="ts">
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import ARecordCard from './components/arecordCard.vue';

import { useRecordsStore } from './stores/records';
import { mapState } from 'pinia'

import { Record } from './models/record';


export default {
  name: 'App',
  components: {
    FontAwesomeIcon,
    ARecordCard,
  },
  setup() {
    const tundraStore = useRecordsStore()
    tundraStore.fetchRecords((_) => { }, (err) => {
      console.log(err);
    })
  },
  data() {
    return {
      allRecords: new Array<Record>(),
      activeRecord: 0,
    };
  },
  methods: {
    getRecordComponent(type: number) {
      const componentMap: { [t: number]: string } = {
        1: "ARecordCard",
        // 2: "CNAMERecordComponent"
        // Add more mappings for other record types as needed
      };
      return componentMap[type] || "div";
    },
    handlePageClick(event: any) {
      let target = event.target;

      // Traverse up the DOM tree and check if any parent has the id starting with "record-"
      while (target !== this.$el) {
        if (target.id && target.id.startsWith('record-')) {

          // Get the record id from the element id
          const recordId = parseInt(target.id.split('-')[1])
          this.activeRecord = recordId

          return
        }
        target = target.parentElement;
      }

      // If we reach here, the user clicked outside of a record card
      this.activeRecord = 0
    },
  },
  computed: {
    ...mapState(useRecordsStore, ['getRecords']),
    getActiveRecord() {
      return this.activeRecord
    },
  },
};

</script>

<template>
  <div class="w-full min-h-max" @click="handlePageClick">
    <div class="container mx-auto max-w-screen-md">
      <!-- Logo -->
      <div class="w-full flex mt-24">
        <img src="./assets/logo-lightmode.svg" class="mx-auto" alt="TundraDNS Logo" />
      </div>

      <!-- Welcome Text -->
      <div class="w-full flex mt-16">
        <p class="text-primary font-sans">
          Welcome to the TundraDNS Portal! Our single-page web application empowers you to effortlessly oversee your
          active domains and their associated records. You currently have the following domain active, you can switch
          between your other domains via the drop down.
        </p>
      </div>

      <!-- Domain Box -->
      <div class="w-full flex mt-8">
        <div class="w-full h-24 flex bg-secondary rounded-lg justify-between">
          <div class="h-full flex ml-5">
            <div class="w-16 h-16 flex my-auto bg-primary-full rounded-lg">
              <font-awesome-icon icon="cloud" class="text-white my-auto mx-auto h-2/5" />
            </div>
            <div class="my-auto ml-8">
              <p class="text-xl  text-primary">
                tundra-dns.io
              </p>
              <div class="flex">
                <div class="border-2 border-success-full rounded-full bg-success-faint">
                  <p class="text-success text-xs px-5 py-0.5">active</p>
                </div>
              </div>
            </div>
          </div>
          <div class="h-full flex">
            <font-awesome-icon icon="chevron-right" class="text-primary my-auto mr-10 cursor-pointer" size="xl" />
          </div>
        </div>
      </div>

      <!-- Domain Records -->
      <div class="w-full flex mt-8">
        <div class="w-full flex flex-col bg-secondary rounded-lg justify-between pt-5">
          <div v-for="record in getRecords" :key="record.id" :id="'record-' + record.id">
            <component :is="getRecordComponent(record.type)" :record="record" :activeRecord="getActiveRecord" />
          </div>

          <div class="w-full">
            <div
              class="w-fill group h-16 justify-center cursor-pointer flex mx-5 mb-5 px-8 rounded-xl border-2 border-primary-full bg-primary-faint border-dashed hover:border-solid hover:bg-primary-full">
              <p class="text-center my-auto text-faint group-hover:text-white ">+ New Record</p>
            </div>
          </div>

        </div>
      </div>

    </div>
  </div></template>
