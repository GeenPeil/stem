import { Pipe, PipeTransform } from '@angular/core';

import { Date } from './date';

/*
 * Convert a Date type to human-readable string (yyyy-mm-dd).
 * Usage:
 *   myDate | dateToString
*/
@Pipe({ name: 'dateToString' })
export class DateToStringPipe implements PipeTransform {
    transform(date: Date): string {
        return date.year + `-` + (date.month < 10 ? '0' : '') + date.month + `-` + (date.day < 10 ? '0' : '') + date.day;
    }
}