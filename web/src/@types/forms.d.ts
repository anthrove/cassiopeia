// Describes a form by what fields it should return
type FormDescriptor = {
    [key:string]: FormFieldDescriptor
}

// Describes a single Form field with optional additional configuration
type ExplicitFormFieldDescriptor = {
    type: "string"
    multiline?: boolean
    min?: number
    max?:number

    default?: any
    label?: string
    required?: boolean
}|{
    type: "number"
    min?: number
    max?: number
    step?: number

    default?: any
    label?: string
    required?: boolean
} | {
    type: "boolean"

    default?: any
    label?: string
    required?: boolean
}| {
    type: "select:single"

    options: any[],
    
    default?: any
    label?: string
    required?: boolean    
}

type FormFieldDescriptor = ExplicitFormFieldDescriptor | StringConstructor | NumberConstructor | BooleanConstructor| string[]