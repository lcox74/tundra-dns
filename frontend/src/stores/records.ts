import { defineStore } from "pinia";
import { RecordType, Record, ARecord } from "../models/record";

export const useRecordsStore = defineStore({
    id: "records",
    state: () => ({
        records: [] as Record[],
    }),
    getters: {
        getRecords(state): Record[]{
            return state.records;
        }
    },
    actions: {
        fetchRecords(sucessCb: (records: Record[]) => void, errorCb: (err: Error) => void) {
            fetch('http://localhost:8053/api/records')
                .then((response) => response.json())
                .then((data) => {

                    this.records = new Array();

                    data.map((data: any) => {
                        if (!data.type) {
                            return
                        }

                        switch (data.type) {
                            case RecordType.A:
                                this.records.push(new ARecord(data));
                                break;
                            default:
                                console.log("Unknown record type: " + data.type);
                        }
                    })

                    sucessCb(this.records);
                })
                .catch((err) => {
                    errorCb(err);
                });
        }
    }
});