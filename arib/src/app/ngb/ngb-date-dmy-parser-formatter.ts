import { NgbDateParserFormatter, NgbDateStruct } from "@ng-bootstrap/ng-bootstrap";

function toInteger(value: any): number {
	return parseInt(`${value}`, 10);
}

function isNumber(value: any): boolean {
	return !isNaN(toInteger(value));
}

function padNumber(value: number) {
	if (isNumber(value)) {
		return `0${value}`.slice(-2);
	} else {
		return "";
	}
}

// NgbDateDMYParserFormatter parses and formats NgbDateStruct from/to dd-mm-yyy format.
export class NgbDateDMYParserFormatter extends NgbDateParserFormatter {
	public parse(value: string): NgbDateStruct {
		if (value) {
			const dateParts = value.trim().split("-");
			if (dateParts.length === 1 && isNumber(dateParts[0])) {
				return { year: toInteger(dateParts[0]), month: null, day: null };
			} else if (dateParts.length === 2 && isNumber(dateParts[0]) && isNumber(dateParts[1])) {
				return { year: toInteger(dateParts[1]), month: toInteger(dateParts[0]), day: null };
			} else if (dateParts.length === 3 && isNumber(dateParts[0]) && isNumber(dateParts[1]) && isNumber(dateParts[2])) {
				return { year: toInteger(dateParts[2]), month: toInteger(dateParts[1]), day: toInteger(dateParts[0]) };
			}
		}
		return null;
	}

	public format(date: NgbDateStruct): string {
		return date ?
			`${isNumber(date.day) ? padNumber(date.day) : ""}-${isNumber(date.month) ? padNumber(date.month) : ""}-${date.year}` :
			"";
	}
}
