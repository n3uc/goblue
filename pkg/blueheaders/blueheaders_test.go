package blueheaders_test

/*
Midas Blue File Header Reading Service

0x000 Start of Header Control Block and Fixed Header
0x100 End of Fixed Header and Start of Variable Header
0x200 End of Variable Header and Header Control Block

... May have 0 or more bytes here, data does not always start right after header

data_start

data_size+data_start

.. May have 0 or more bytes here,  Extended header must start on 512 Byte boundary

512*ext_start    Start of Extended Header
ext_size + 512*ext_start    End of Extended Header

Note:  Extended Header may or may not come after data section

Fixed Header
------------
0 	Version			char[4]
4	head_rep		char[4]
8	data rep		char[4]
12	detached		int_4
16	protected		int_4
20	pipe			int_4
24	ext_start		int_4
28	ext_size		int_4
32	data_start		real_8
40	data_size		real_8
48 	type			int_4
52	format			char[2]
54	flagmask		int_2
56	timecode		real_8
64	inlet			int_2
66	outlets			int_2
68	outmask			int_4
72	pipeloc			int_4
76	pipesize		int_4
80	in_byte			real_8
88	out_byte		real_8
96	outbytes		real_8[8]
160	keywords		char[92]
256	adjunct			char[256]
*/

// func TestFull(t *testing.T) {
// 	header, err := blueheaders.New("sample.tmp")
// 	if err != nil {
// 		t.Fatalf("Could not load header\n%s\n", err)
// 	}

// 	fmt.Printf("Loaded full header\n%+v\n", header)

// 	if header.ExtendedHeaders["KEYWORD1"] != "ONE" {
// 		t.Fatal("KEYWORD1 is incorrect")
// 	}
// }
