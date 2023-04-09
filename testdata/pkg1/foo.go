package pkg1

//go:generate go run ../../ -package ./ -outdir ../mock1_gen1 -revision 1 Foo
//go:generate go run ../../ -package ./ -outdir ../mock1_gen2 -revision 2 Foo
//go:generate go run ../../ -package ./ -outdir ../mock1_gen3 -revision 3 Foo

type Foo struct {
}

func (*Foo) Hello(name string) error {}
