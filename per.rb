#!/usr/bin/env ruby
# encoding: utf-8

per = Proc.new { |len, num| num / len.to_f * 100.floor }

STDIN.readlines.each_with_index do |s, i|
  s.chomp!
  len = s.length

  printf("- %5d -\n[%-7s]: %s\n", i+1, 'input', s)

  l = s.count("a-z")
  u = s.count("A-Z")
  n = s.count("0-9")
  o = len-l-u-n

  format = "[%-7s]: [%4d]\n" + "[%-7s]: [%4d, %3d%%]\n"*4

  printf(format,
         'length', len,
         'lower', l, per.call(len, l),
         'upper', u, per.call(len, u),
         'number', n, per.call(len, n),
         'others', o, per.call(len, o))
  puts
end
