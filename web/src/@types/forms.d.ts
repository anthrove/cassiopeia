// Describes a form by what fields it should return
type FormDescriptor = {
    [key:string]: FormFieldDescriptor
}

type LabelerFunction = (string)=>string
type CommonFieldDescriptorOptions = {
    default: any
    label: string
    required: boolean,
    readonly: boolean,
    conditional? = (formState:any)=>boolean
}

// Describes a single Form field with optional additional configuration

type ExplicitFormFieldDescriptor = Partial<CommonFieldDescriptorOptions>&({
    type: "string"
    multiline?: boolean
    min?: number
    max?:number
}|{
    type: "number"
    min?: number
    max?: number
    step?: number

} | {
    type: "boolean"
}| {
    type: "date"
} | {
    type: "select:single"
    options: any[],
    label: string
    labeler?: LabelerFunction
})

type FormFieldDescriptor = ExplicitFormFieldDescriptor | StringConstructor | NumberConstructor | BooleanConstructor| string[] | DateConstructor