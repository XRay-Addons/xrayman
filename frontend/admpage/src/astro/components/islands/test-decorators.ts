function dec(value: any, context: ClassAccessorDecoratorContext) {
  console.log(context.name);
}

export class Test {
  @dec accessor x = 1;
}
