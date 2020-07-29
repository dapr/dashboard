import { Pipe, PipeTransform } from '@angular/core';
import { Log } from 'src/app/types/types';

@Pipe({
    name: 'sort'
})
export class SortPipe implements PipeTransform {
    transform(logs: Log[], order: string): Log[] {
        console.log(order);
        if (!order) {
            return logs;
        }
        const comparator = (a: Log, b: Log) => {
            const aTime = (a.timestamp !== "") ? 0 : Date.parse(a.timestamp);
            const bTime = (b.timestamp !== "") ? 0 : Date.parse(b.timestamp);
            return aTime - bTime;
        }
        if (order === 'asc') {
            return logs.sort(comparator);
        } else {
            return logs.sort(comparator).reverse();
        }
    }
}
