<script lang="ts">
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';

import { ARecord } from '../models/record';

export default {
    name: 'ARecordCard',
    components: {
        FontAwesomeIcon,
    },
    props: {
        record: {
            type: ARecord,
            required: true,
        },
        activeRecord:  {
            type: Number,
            required: true,
        },
    },
    data() {
        return {
            timeSinceLastSeen: 0,
            timeSinceInterval: 0,

            selected: false,
        }
    },
    mounted() {
        this.timeSinceInterval = setInterval(() => {
            this.setTimeSinceLastSeen()
        }, 1000)
    },
    beforeDestroy() {
        clearInterval(this.timeSinceInterval)
    },
    watch: {
        activeRecord: {
            handler: function (val: number, _) {
                if (val == this.record.id) {
                    this.selected = true
                } else {
                    this.selected = false
                }
            },
            immediate: true,
        },
    },
    methods: {
        setTimeSinceLastSeen() {
            this.timeSinceLastSeen = Math.floor((Date.now() - this.record.lastSeen.getTime()) / 1000)
        },
    },
    computed: {
        getSelected() {
            return this.selected
        },
        getTimeSinceLastSeen() {
            return this.timeSinceLastSeen 
        },
    }
}

</script>

<template>
    <div class="w-fill min-h-24 mx-5 mb-5 px-8 rounded-xl border-2 border-gray-200 hover:border-2 hover:border-solid hover:border-accent"
        :class="[getSelected ? 'border-solid border-primary-full' : 'border-dashed cursor-pointer']" >
        <div class="flex w-full h-16 ">
            <div class="grid w-full grid-cols-10 font-mono my-auto">
                <div>
                    <p class="text-failure font-bold text-center text-3xl">A</p>
                </div>
                <div class="col-span-1"></div>
                <div class="col-span-2 flex">
                    <p class="text-primary my-auto font-medium text-center">{{ record.subdomain }}</p>
                </div>
                <div class="col-span-2"></div>
                <div class="flex justify-end col-span-4 ">
                    <p class="text-primary my-auto font-medium text-center">{{ record.address }}</p>
                    <font-awesome-icon icon="circle" class="my-auto text-success ml-2" size="2xs" />
                </div>
            </div>
        </div>
        <div v-if="getSelected">
            <div class="border border-primary-full"></div>
            <div class="flex justify-between font-mono text-primary my-4">
                <div class="flex">
                    <p class="font-semibold">FQDN:</p>
                    <p class="ml-1">{{ record.subdomain }}.{{ record.domain }}</p>
                </div>
                <div class="flex">
                    <p class="font-semibold">TTL:</p>
                    <p class="ml-1">{{ record.ttl }}s</p>
                </div>
                <div class="flex">
                    <p class="font-semibold">Last Seen:</p>
                    <p class="ml-1">{{ getTimeSinceLastSeen }}s</p>
                </div>
            </div>
            <div class="flex justify-end mb-4">
                <div class="bg-primary-full border-2 border-primary-full hover:border-white text-white font-mono cursor-pointer rounded-md w-32 p-1 text-sm text-center mr-4">
                    <p class="select-none">Edit Node</p>
                </div>
                <div class="bg-primary-full border-2 border-primary-full hover:border-white text-white font-mono cursor-pointer rounded-md w-32 p-1 text-sm text-center">
                    <p class="select-none">Delete Node</p>
                </div>
            </div>
        </div>
    </div>
</template>